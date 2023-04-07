package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "1"
	dbname   = "mydb"
)

type Products struct {
	Id      int
	Model   string
	Company string
	Price   int
}

var database *sql.DB

func main() {
	dns := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/create", createHandler)
	router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)
	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("Get")
	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server is starting at 127.0.0.1:4000")
	err = http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
