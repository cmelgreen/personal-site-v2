package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"personal-site-v2/backend/server/database"

	"github.com/go-chi/chi"
	"golang.org/x/crypto/acme/autocert"
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
	db         *database.Database
	mux        chi.Router
	middleware chi.Middlewares
	log        *log.Logger
}

// NewServer returns new server with default log, mux, and database
func newServer(ctx context.Context) *Server {
	s := Server{
		log: log.New(logOut, logPrefix, logFlags),
		mux: chi.NewRouter(),
		db:  &database.Database{},
	}

	return &s
}

func (s *Server) printRoutes() {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		s.log.Printf("%s %s \n", method, route)
		return nil
	}
	
	if err := chi.Walk(s.mux, walkFunc); err != nil {
		s.log.Printf("Logging err: %s\n", err.Error())
	}
}

func (s *Server) serve(port string) {
	s.log.Fatal(http.ListenAndServe(port, s.mux))
}

func (s *Server) serveHTTPS(port string) {
	cert := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.cmelgreen.com"),
		Cache:  autocert.DirCache("../."),
	}

	httpsMux := &http.Server{
		Addr:    ":443",
		Handler: s.mux,
		TLSConfig: &tls.Config{
			GetCertificate: cert.GetCertificate,
		},
	}

	go http.ListenAndServe(":80", cert.HTTPHandler(nil))
	httpsMux.ListenAndServeTLS("", "")

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

func status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("okay"))
}