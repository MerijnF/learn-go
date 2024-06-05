package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./todo.db"
const dbDriver = "sqlite3"

var db *sql.DB

type Todo struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Completed bool 	`json:"completed"`
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func handleGet(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	stmt, err := db.Prepare("SELECT id, title, completed FROM todos")
	check(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	check(err)

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.Id, &todo.Title, &todo.Completed)
		check(err)
		todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)
}

func handleGetWithId(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	pathId := req.PathValue("id")

	stmt, err := db.Prepare("SELECT id, title, completed FROM todos WHERE id = ?")
	check(err)
	defer stmt.Close()

	row, err := stmt.Query(pathId)
	check(err)

	var todo Todo
	for row.Next() {
		err = row.Scan(&todo.Id, &todo.Title, &todo.Completed)
		check(err)
	}
	json.NewEncoder(w).Encode(todo)
}

func handlePut(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	pathId := req.PathValue("id")

	var todo Todo
	err := json.NewDecoder(req.Body).Decode(&todo)
	check(err)

	stmt, err := db.Prepare("UPDATE todos SET title = ?, completed = ? WHERE id = ?")
	check(err)
	defer stmt.Close()
	_, err = stmt.Exec(todo.Title, todo.Completed, pathId)
	check(err)

	todo.Id, err = strconv.Atoi(pathId)
	check(err)

	json.NewEncoder(w).Encode(todo)
}

func handlePost(w http.ResponseWriter, req *http.Request) {
	setHeaders(w)

	var todo Todo
	err := json.NewDecoder(req.Body).Decode(&todo)

	stmt, err := db.Prepare("INSERT INTO todos(title, completed) VALUES(?, ?)")
	check(err)
	defer stmt.Close()
	stmtResult, err := stmt.Exec(todo.Title, todo.Completed)

	id, err := stmtResult.LastInsertId()
	check(err)

	todo.Id = int(id)
	
	json.NewEncoder(w).Encode(todo)
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func handleDelete(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	pathId := req.PathValue("id")

	stmt, err := db.Prepare("DELETE FROM todos WHERE id = ?")
	check(err)
	defer stmt.Close()
	_, err = stmt.Exec(pathId)
	check(err)
}

func main() {
	arg := os.Args[1:]
	if(len(arg) > 0 && arg[0] == "clear") {
		os.Remove(dbPath)
	}

	dbConnection, err := sql.Open(dbDriver, dbPath)
	check(err)
	db = dbConnection
	defer db.Close()

	fmt.Println("Creating table (if it doesn't exist)")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, title TEXT, completed BOOLEAN)")
	check(err)

	fmt.Println("Serving on http://localhost:8090")
	http.HandleFunc("GET /todo" , handleGet)
	http.HandleFunc("POST /todo", handlePost)
	http.HandleFunc("GET /todo/{id}" , handleGetWithId)
	http.HandleFunc("PUT /todo/{id}" , handlePut)
	http.HandleFunc("DELETE /todo/{id}" , handleDelete)
	http.HandleFunc("OPTIONS /todo", func(w http.ResponseWriter, r *http.Request) {setHeaders(w)})
	http.HandleFunc("OPTIONS /todo/{id}", func(w http.ResponseWriter, r *http.Request) {setHeaders(w)})
	

	http.ListenAndServe(":8090", nil)
}