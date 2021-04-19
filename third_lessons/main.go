package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts = []Article{}
var showPost = Article{}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golangdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err.Error())
	}

	posts = []Article{}
	for res.Next() {
		var article Article
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			panic(err.Error())
		}

		posts = append(posts, article)
	}

	defer res.Close()

	tmpl.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")

	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Not all fields are filled!")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golangdb")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`Title`, `Anons`, `FullText`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err.Error())
		}
		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/golangdb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = %s", vars["id"]))
	if err != nil {
		panic(err.Error())
	}

	showPost = Article{}
	for res.Next() {
		var article Article
		err = res.Scan(&article.Id, &article.Title, &article.Anons, &article.FullText)
		if err != nil {
			panic(err.Error())
		}

		showPost = article
	}

	defer res.Close()

	tmpl.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {
	r := mux.NewRouter()
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/create", create).Methods("GET")
	r.HandleFunc("/save_article", save_article).Methods("POST")
	r.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", r)
}

func main() {
	handleFunc()
}
