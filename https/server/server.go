package https

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"

	"github.com/zayaanra/thunderspeak/api"

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

	rooms map[*api.Room]bool
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),

		Addr:   "localhost:8080",
		Domain: "localhost",

		Dir: "../https/frontend/src",

		rooms: make(map[*api.Room]bool),
	}

	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	s.router.HandleFunc("/", s.handleIndex)
	s.router.HandleFunc("/ws", s.handleWS)

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
		log.Printf("POST - %s\n", r.URL.Path)
		switch r.URL.Path {
		case "/api/createRoom":
			s.createRoom()
		}
	} else if r.Method == http.MethodGet {
		log.Printf("GET - %s\n", r.URL.Path)
	}

	s.router.ServeHTTP(w, r)

}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.Dir+"/index.html")
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	api.ServeWS(w, r)
}

func (s *Server) createRoom() *api.Room {
	roomName := rand.Intn(1000000)
	room := api.NewRoom(fmt.Sprintf("%d", roomName))

	log.Printf("Created room - %v\n", room)

	return room
}
