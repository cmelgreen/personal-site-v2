package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"PersonalSite/backend/database"

	"github.com/rs/cors"
	"github.com/julienschmidt/httprouter"
)

const (
	heartbeatTime = 10
)

var (
	logOut    = os.Stdout
	logPrefix = log.Prefix()
	logFlags  = log.Flags()
)

// Server struct for storing database, mux, and logger
type Server struct {
	db  *database.Database
	mux *httprouter.Router
	log *log.Logger
}

// NewServer returns new server with default log, mux, and database
func newServer(ctx context.Context) *Server {
	s := Server{
		log: log.New(logOut, logPrefix, logFlags),
		mux: httprouter.New(),
		db:  &database.Database{},
	}

	return &s
}

func (s *Server) serve(port string) {
	s.log.Fatal(http.ListenAndServe(port, s.mux))
}

func (s *Server) serveCORSEnabled(port string) {
	muxCORS := cors.AllowAll().Handler(s.mux)
	s.log.Fatal(http.ListenAndServe(port, muxCORS))
}

// NewDBConnection creates a new connection to a database for a server
func (s *Server) newDBConnection(ctx context.Context, dbConfig database.DBConfig) {
	var err error

	// FIX NULL ERRORS
	s.db, err = database.ConnectToDB(ctx, dbConfig)
	if err != nil {
		s.log.Println(err)
	}

	//s.maintainDBConnection(ctx, dbConfig)
}

func (s *Server) maintainDBConnection(ctx context.Context, dbConfig database.DBConfig) {
	go func() {
		var err error
		for {
			if s.db.Connected(ctx) != true {
				s.db, err = database.ConnectToDB(ctx, dbConfig)
				if err != nil {
					s.log.Println("Error maintaining connection: ", err)
				}
			}
			time.Sleep(heartbeatTime * time.Second)
		}
	}()
}
