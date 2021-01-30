package main

import (
	"context"
	"log"
	"os"
	"time"

	"PersonalSite/backend/database"

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
