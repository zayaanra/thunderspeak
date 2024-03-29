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
	s.router.HandleFunc("/room/{room_code:[0-9]+}", s.handleRoom)

	s.router.HandleFunc("/ws", s.handleWS)

	fs := http.FileServer(http.Dir(s.Dir))
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return s
}

// Opens a new server
func (s *Server) Open() (err error) {
	s.ln, err = net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	go s.server.Serve(s.ln)

	return nil
}

// Shuts down the server
func (s *Server) Close() error {
	return s.ln.Close()
}

// Listens for API HTTP requests and responds accordingly
func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

	} else if r.Method == http.MethodGet {
		switch r.URL.Path {
		case "/api/createRoom":
			room := s.createRoom()
			http.Redirect(w, r, "/room/"+room.Name(), http.StatusSeeOther)
			return
		}
	}

	s.router.ServeHTTP(w, r)

}

// Handles serving up index.html as the homepage
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.Dir+"/index.html")
}

// Handles serving up room.html to represent a room
func (s *Server) handleRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room_code := vars["room_code"]

	if !s.roomExists(room_code) {
		http.NotFound(w, r)
	} else {
		http.ServeFile(w, r, s.Dir+"/room.html")
	}

}

// Handles API request to switch to WS protocol
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	api.ServeWS(w, r)
}

// Handles API request to create a room
func (s *Server) createRoom() *api.Room {
	// Generate random room code and create room
	roomName := rand.Intn(1000000)
	room := api.NewRoom(fmt.Sprintf("%d", roomName))

	// Mark room as false to indicate that it is empty
	s.rooms[room] = false

	log.Printf("Created room - %v\n", room)

	return room
}

func (s *Server) roomExists(room_code string) bool {
	for room := range s.rooms {
		if room.Name() == fmt.Sprintf("%v", room_code) {
			return true
		}
	}
	return false
}
