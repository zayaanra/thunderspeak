package https

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	// Basic server needs. A listener, server, and router.
	ln     net.Listener
	server *http.Server
	router *mux.Router

	// Bind address & domain for the server's listener.
	Addr   string
	Domain string

	// The server's directory to serve HTML files from.
	Dir string
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),

		Addr:   "localhost:8080",
		Domain: "localhost",

		Dir: "../https/frontend/src",
	}

	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// Set up handling of invalid routes.
	s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)

	s.router.HandleFunc("/", s.handleIndex)

	fs := http.FileServer(http.Dir(s.Dir))
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return s
}

func (s *Server) Open() (err error) {
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	return nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("POST")
	} else if r.Method == http.MethodGet {
		log.Println("GET")
	}

	s.router.ServeHTTP(w, r)

}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.Dir+"/index.html")
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	log.Println("NOT FOUND")
	http.NotFound(w, r)
}
