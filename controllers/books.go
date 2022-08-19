package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type BooksController struct {
	DB         *sql.DB
	Router     *mux.Router
	PathPrefix string
}

// BOOKS
func ReadPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]

	fmt.Fprintf(w, "You've requested the book: %s, page %s\n", title, page)
}

func AllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've get all books\n")
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've created the book: %s\n", title)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've readed the book: %s\n", title)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've upated the book: %s\n", title)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've deleted the book: %s\n", title)
}

// Register routes function
func (c BooksController) Register() {
	bookrouter := c.Router.PathPrefix(c.PathPrefix).Subrouter()
	bookrouter.HandleFunc("/", AllBooks)
	bookrouter.HandleFunc("/{title}", GetBook).Methods("GET")
	bookrouter.HandleFunc("/{title}", CreateBook).Methods("POST")
	bookrouter.HandleFunc("/{title}", UpdateBook).Methods("PUT")
	bookrouter.HandleFunc("/{title}", DeleteBook).Methods("DELETE")
	bookrouter.HandleFunc("/{title}/page/{page}", ReadPage).Methods("GET")
}
