package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func productHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	cat := vars["category"]
	response := fmt.Sprintf("Product category=%s id=%s", cat, id)
	fmt.Fprint(w, response)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from products")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []Products{}

	for rows.Next() {
		p := Products{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			log.Println(err)
		}
		products = append(products, p)
	}
	tpml, err := template.ParseFiles("./html/index.html")
	tpml.Execute(w, products)
}
func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Panicln(err)
		}
		model := r.FormValue("model")
		company := r.FormValue("company")
		price := r.FormValue("price")

		_, err = database.Exec("insert into products (model,company,price) values($1,$2,$3)", model, company, price)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "./html/create.html")
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec("delete from products where id=$1", id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)

}

func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from products where id=$1", id)
	prod := Products{}
	err := row.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("./html/edit.html")
		tmpl.Execute(w, prod)
	}
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	model := r.FormValue("model")
	company := r.FormValue("company")
	price := r.FormValue("price")

	_, err = database.Exec("update products set model=$1,company=$2,price=$3 where id=$4", model, company, price, id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}
