package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

type User struct {
	Id   int
	Name string
	Age  int
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

		log.Println("User")
		log.Println(dbUser)

		users = append(users, dbUser)
	}

	return users
}

func seedUsers(db *sql.DB) {
	users := []User{
		{1, "Brian", 19},
		{2, "Thomas", 32},
		{3, "Tonny", 99},
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS user (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255), age INT)")
	if nil != err {
		log.Fatal(err)
	}

	log.Println("Inserting...")
	for _, user := range users {
		_, err = db.Exec("INSERT INTO user (name, age) values (?, ?)", user.Name, user.Age)
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
	r.Handle("/", http.RedirectHandler("/hello.html", 301))

	r.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		w.Header()["content-type"] = []string{"application/json"}
		fmt.Fprint(w, toJson(fetchUsers(db)))
	})

	r.HandleFunc("/data/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)

		if id <= 0 {
			http.NotFound(w, req)
			return
		}

		row := db.QueryRow("SELECT name, age FROM hello_go.user where id = ?", id)

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

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
