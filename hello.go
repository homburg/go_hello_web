package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lann/squirrel"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

type User struct {
	Id   int
	Name string
	Age  int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func pushUsers(db *sql.DB, conns []*websocket.Conn) {
	users := fetchUsers(db)
	for _, conn := range conns {
		if nil == users {
			users = []User{}
		}
		websocket.WriteJSON(conn, users)
	}
}

func fetchUsers(db *sql.DB) []User {
	dbUsers, err := db.Query("SELECT * FROM user")
	if nil != err {
		log.Fatal(err)
	}
	defer dbUsers.Close()

	var users []User
	for dbUsers.Next() {
		var id int
		var name string
		var age int
		if dbUsers.Scan(&id, &name, &age); nil != err {
			log.Fatal(err)
		}

		dbUser := User{id, name, age}

		users = append(users, dbUser)
	}

	return users
}

func seedUsers(db *sql.DB) {
	users := []User{
		{0, "Brian", 19},
		{0, "Thomas", 32},
		{0, "Tonny", 99},
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), age INT)")
	if nil != err {
		log.Fatal(err)
	}

	log.Println("Inserting...")
	for _, user := range users {
		_, err = squirrel.Insert("user").Columns("name", "age").Values(user.Name, user.Age).RunWith(db).Exec()
		if nil != err {
			log.Fatal(err)
		}
	}
}

func toJson(data interface{}) string {
	json, _ := json.MarshalIndent(data, "", "  ")
	return string(json)
}

func main() {
	var conns []*websocket.Conn

	log.Println("Starting...")
	runtime.GOMAXPROCS(runtime.NumCPU())
	var err error

	// DB conn
	db, err := sql.Open("mysql", "hello_go@/hello_go")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	seedUsers(db)

	r := mux.NewRouter()

	r.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}

		conns = append(conns, conn)
	})

	r.Methods("GET", "HEAD").Path("/data").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		users := fetchUsers(db)
		w.Header()["content-type"] = []string{"application/json"}
		w.Header()["x-count"] = []string{strconv.FormatInt(int64(len(users)), 10)}
		if len(users) == 0 {
			w.Write([]byte("[]"))
		} else {
			fmt.Fprint(w, toJson(users))
		}
	})

	r.Methods("GET", "HEAD").Path("/data/{id}").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		if id <= 0 {
			http.NotFound(w, req)
			return
		}

		row := squirrel.Select("name", "age").From("user").Where(squirrel.Eq{"id": id}).RunWith(db).QueryRow()

		var name string
		var age int
		if err = row.Scan(&name, &age); err != nil {
			log.Println(err)
			http.NotFound(w, req)
			return
		}

		w.Header()["content-type"] = []string{"application/json"}
		fmt.Fprint(w, toJson(User{int(id), name, age}))
	})

	r.Methods("DELETE").Path("/data").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Deleting...")
		_, err := squirrel.Delete("user").RunWith(db).Exec()
		if nil != err {
			log.Println(err)
		}
		pushUsers(db, conns)
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprint(w, "[]")
	})

	r.Methods("POST").Path("/data").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Seeding...")
		seedUsers(db)
		pushUsers(db, conns)
		w.WriteHeader(http.StatusCreated)
	})

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(rice.MustFindBox("public").HTTPBox()),
	)
	n.UseHandler(r)
	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":3000"
	}
	n.Run(listen)
}
