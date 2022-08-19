package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type AuthorsController struct {
	DB         *sql.DB
	Router     *mux.Router
	PathPrefix string
}

// AUTHORS
func (c AuthorsController) AllAuthors(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've get all authors\n")
}

func (c AuthorsController) GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Fprintf(w, "You've get the author: %s\n", name)
}

// Register routes function
func (c AuthorsController) Register() {
	authorrouter := c.Router.PathPrefix(c.PathPrefix).Subrouter()
	authorrouter.HandleFunc("/", c.AllAuthors)
	authorrouter.HandleFunc("/{name}", c.GetAuthor).Methods("GET")
}
