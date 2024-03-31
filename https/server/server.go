package https

import (
	"encoding/json"
	"fmt"
	"io"
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

type Username struct {
	Username string `json:"username"`
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
		switch r.URL.Path {
		case "/api/createRoom":
			room := s.createRoom(w, r)
			http.Redirect(w, r, "/room/"+room.Name(), http.StatusSeeOther)
			return
		case "/api/joinRoom":
			log.Println("Attempting to join room")
		}
		// TODO: Process API request for joining a room
	} //else if r.Method == http.MethodGet {

	// }

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

	if _, exists := s.roomExists(room_code); !exists {
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
func (s *Server) createRoom(w http.ResponseWriter, r *http.Request) *api.Room {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var username Username
	err = json.Unmarshal(body, &username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Generate random room code and create room
	roomName := rand.Intn(1000000)
	room := api.NewRoom(fmt.Sprintf("%d", roomName))

	// Client that created the room joins the room. Mark room as true to indicate it is not empty.
	s.joinRoom(username.Username, room.Name())
	s.rooms[room] = true

	return room
}

// Handles joining of a room. The given user is added to the room with the given room code.
func (s *Server) joinRoom(username, room_code string) {
	client := api.NewClient(username, nil)
	room, _ := s.roomExists(room_code)

	// Register the client into the room.
	room.Register(client)
}

func (s *Server) roomExists(room_code string) (*api.Room, bool) {
	for room := range s.rooms {
		if room.Name() == fmt.Sprintf("%v", room_code) {
			return room, true
		}
	}
	return nil, false
}
