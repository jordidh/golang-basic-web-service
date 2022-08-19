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
func (c BooksController) ReadPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	page := vars["page"]

	fmt.Fprintf(w, "You've requested the book: %s, page %s\n", title, page)
}

func (c BooksController) AllBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've get all books\n")
}

func (c BooksController) CreateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've created the book: %s\n", title)
}

func (c BooksController) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've readed the book: %s\n", title)
}

func (c BooksController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've upated the book: %s\n", title)
}

func (c BooksController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "You've deleted the book: %s\n", title)
}

// Register routes function
func (c BooksController) Register() {
	bookrouter := c.Router.PathPrefix(c.PathPrefix).Subrouter()
	bookrouter.HandleFunc("/", c.AllBooks)
	bookrouter.HandleFunc("/{title}", c.GetBook).Methods("GET")
	bookrouter.HandleFunc("/{title}", c.CreateBook).Methods("POST")
	bookrouter.HandleFunc("/{title}", c.UpdateBook).Methods("PUT")
	bookrouter.HandleFunc("/{title}", c.DeleteBook).Methods("DELETE")
	bookrouter.HandleFunc("/{title}/page/{page}", c.ReadPage).Methods("GET")
}
