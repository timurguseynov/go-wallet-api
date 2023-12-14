// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// This program provides a sample web service that implements a
// RESTFul CRUD API against a MongoDB database.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/timurguseynov/go-wallet-api/internal/platform/db"

	"github.com/timurguseynov/go-wallet-api/cmd/apid/handlers"
)

const (
	host = "localhost"
	port = "3000"
)

// init is called before main. We are using init to customize logging output.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

// main is the entry point for the application.
func main() {
	log.Println("main : Started")

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB:")
	dbConn, err := db.NewDB()
	if err != nil {
		log.Fatal("main : couldn't connect to database", err)
	} else {
		log.Println("main : DB captured successfully")
	}

	server := http.Server{
		Addr:    host + ":" + port,
		Handler: handlers.API(dbConn),
	}

	// We want to report the listener is closed.
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the listener.
	go func() {
		log.Printf("startup : Listening %s", host)
		log.Printf("shutdown : Listener closed : %v", server.ListenAndServe())
		wg.Done()
	}()

	// Listen for an interrupt signal from the OS.
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	// Wait for a signal to shutdown.
	<-osSignals

	// Create a context to attempt a graceful 5 second shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	// Attempt the graceful shutdown by closing the listener and
	// completing all inflight requests.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown : Graceful shutdown did not complete in %v : %v", time.Duration(5*time.Second), err)

		// Looks like we timedout on the graceful shutdown. Kill it hard.
		if err := server.Close(); err != nil {
			log.Printf("shutdown : Error killing server : %v", err)
		}
	}

	// Wait for the listener to report it is closed.
	wg.Wait()
	log.Println("main : Completed")
}
