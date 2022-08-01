package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const SERVER_IP = "0.0.0.0"

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

// AUTHORS
func AllAuthors(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You've get all authors\n")
}

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	fmt.Fprintf(w, "You've get the author: %s\n", name)
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")

	// Paràmetres que es poden passar a la crida amb el seu valor per defecte
	var env string
	var port int
	var printRoutes bool
	var help bool

	flag.StringVar(&env, "env", "development", "current environment")
	flag.IntVar(&port, "port", 3000, "port number")
	flag.BoolVar(&printRoutes, "print-routes", false, "display all server routes")
	flag.BoolVar(&help, "help", false, "display all server argument available")

	flag.Parse()

	if help {
		// Imprimim els paràmetres per defecte
		fmt.Println()
		fmt.Println("basic-web-server <argument> ...")
		fmt.Println("All arguments:")
		flag.PrintDefaults()
		fmt.Println()

		os.Exit(0)
	}

	// Mostrem el valor dels paràmetres
	fmt.Println("Argument values:")
	fmt.Println("\tenv:", env)
	fmt.Println("\tport:", port)
	fmt.Println("\tprint-routes:", printRoutes)

	// Declarem les rutes
	router := mux.NewRouter()

	bookrouter := router.PathPrefix("/books").Subrouter()
	bookrouter.HandleFunc("/", AllBooks)
	bookrouter.HandleFunc("/{title}", GetBook).Methods("GET")
	bookrouter.HandleFunc("/{title}", CreateBook).Methods("POST")
	bookrouter.HandleFunc("/{title}", UpdateBook).Methods("PUT")
	bookrouter.HandleFunc("/{title}", DeleteBook).Methods("DELETE")
	bookrouter.HandleFunc("/{title}/page/{page}", ReadPage).Methods("GET")

	authorrouter := router.PathPrefix("/authors").Subrouter()
	authorrouter.HandleFunc("/", AllAuthors)
	authorrouter.HandleFunc("/{name}", GetAuthor).Methods("GET")

	// Mostrem totes les
	if printRoutes {
		fmt.Println("Server routes dump:")
		err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()

			if err == nil {
				fmt.Print(" ROUTE:", pathTemplate)
			}
			pathRegexp, err := route.GetPathRegexp()
			if err == nil {
				fmt.Print(" Path regexp:", pathRegexp)
			}
			queriesTemplates, err := route.GetQueriesTemplates()
			if err == nil {
				fmt.Print(" Queries templates:", strings.Join(queriesTemplates, ","))
			}
			queriesRegexps, err := route.GetQueriesRegexp()

			if err == nil {
				fmt.Print(" Queries regexps:", strings.Join(queriesRegexps, ","))
			}
			methods, err := route.GetMethods()
			if err == nil {
				fmt.Print(" Methods:", strings.Join(methods, ","))
			}
			fmt.Println()
			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}

	// Creem el servidor
	srv := &http.Server{
		Addr: SERVER_IP + ":" + fmt.Sprint(port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	fmt.Println()
	log.Println("Server listening at", SERVER_IP+":"+fmt.Sprint(port))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
