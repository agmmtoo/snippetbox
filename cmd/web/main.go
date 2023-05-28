package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// config
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// app instance with di
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
