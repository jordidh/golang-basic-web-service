package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	books "github.com/jordidh/basicwebserver/routes"
	publisher "github.com/jordidh/basicwebserver/routes"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// Estructura per importar el fitxer JSON amb la configuració
type Configuration struct {
	Database struct {
		Driver   string
		Server   string
		Port     int
		User     string
		Password string
		Database string
	}
}

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

		// To exit a program without an error, you can set the exit code of the program to 0
		// Any other numerical value between 1 and 125 (golang) shows the program encountered an error.
		os.Exit(0)
	}

	// Mostrem el valor dels paràmetres
	fmt.Println("Argument values:")
	fmt.Println("\tenv:", env)
	fmt.Println("\tport:", port)
	fmt.Println("\tprint-routes:", printRoutes)

	// Carreguem el fitxer amb la configuració
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Fatal("can't open config file: ", err)
	}
	// Tanquem la connexió abans de finalitzar el programa
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("can't decode config JSON: ", err)
	}
	log.Println("Config file:", Config)

	// Connectem amb la BD
	// parseTime: parseTime=true changes the output type of DATE and DATETIME values to time.Time instead of []byte / string The date or datetime like 0000-00-00 00:00:00 is converted into zero value of time.Time.
	db, err := sql.Open(Config.Database.Driver,
		Config.Database.User+":"+Config.Database.Password+"@("+Config.Database.Server+":"+fmt.Sprint(Config.Database.Port)+")/"+Config.Database.Database+"?parseTime=true")
	if err != nil {
		log.Fatal(err) // Fatal is equivalent to Print() followed by a call to os.Exit(1).
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to dabatase", Config.Database.Database)

	// Tanquem la connexió abans de finalitzar el programa
	defer db.Close()

	var version string
	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err2 != nil {
		log.Fatal(err2)
	}
	log.Println("Database version", version)

	/*
		// Fem una consulta a la BD
		res, err := db.Query("SELECT * FROM cities")

		// Tanquem la connexió a BD
		defer res.Close()

		// Mostrem les dades
		if err != nil {
			log.Fatal(err)
		}

		for res.Next() {
			var city City
			err := res.Scan(&city.Id, &city.Name, &city.Population)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%v\n", city)
		}
	*/

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

	pt1 := publisher.Point{X: 2, Y: 3}
	fmt.Println(pt1)

	pt2 := books.Point{X: 2, Y: 3}
	fmt.Println(pt2)

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
