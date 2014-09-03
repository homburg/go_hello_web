package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
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

func toJson(data interface{}) string {
	json, _ := json.MarshalIndent(data, "", "  ")
	return string(json)
}

func main() {
	log.Println("Starting...")

	runtime.GOMAXPROCS(runtime.NumCPU())

	users := []User{
		{1, "Brian", 19},
		{2, "Thomas", 32},
		{3, "Tonny", 99},
	}

	// DB conn
	db, err := sql.Open("sqlite3", "./hello.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, name VARCHAR(255), age INTEGER)")
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

	dbUsers, err := db.Query("SELECT * FROM user")
	if nil != err {
		log.Fatal(err)
	}
	defer dbUsers.Close()

	log.Println("dbUsers")
	log.Println(dbUsers)

	users = make([]User, 0)
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

	r := mux.NewRouter()
	r.Handle("/", http.RedirectHandler("/hello.html", 301))

	r.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, toJson(users))
	})

	r.HandleFunc("/data/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id, _ := strconv.ParseInt(vars["id"], 10, 0)
		id = id - 1

		if id < 0 || int(id) > (len(users)-1) {
			http.NotFound(w, req)
		} else {
			fmt.Fprintf(w, toJson(users[id]))
		}
	})

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
