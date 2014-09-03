package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	users := []User{
		{1, "Brian", 19},
		{2, "Thomas", 32},
		{3, "Tonny", 99},
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
