package api

type Room struct {
	name       string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan string
}

func NewRoom(name string) *Room {
	return &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan string),
	}
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Register(c *Client) {
	r.clients[c] = true
}

func (r *Room) Unregister(c *Client) {
	delete(r.clients, c)
}

func (r *Room) Run() {
	for {
		select {
		// Register the client if we received a message from the register channel
		case client := <-r.register:
			r.clients[client] = true
		// Unregister the client if we received a message from the unregister channel
		case client := <-r.unregister:
			delete(r.clients, client)
			// TODO: Broadcast the message to all clients in the room
			// case msg := <-r.broadcast:
			// for c := range r.clients {
			// 	c.Send(client)
			// }
		}
	}
}
