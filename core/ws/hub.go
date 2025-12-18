package ws

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

// Hub maintains the set of active clients
type Hub struct {
	// Registered clients needed to send to specific user
	// Map Link UserAccessID -> []Connections (One user might have multiple tabs)
	Clients       map[int][]*websocket.Conn
	Register      chan *RegisterInfo
	Unregister    chan *RegisterInfo
	Broadcast     chan []byte
	mutex         sync.Mutex
	OnUserOnline  func(userID int)
	OnUserOffline func(userID int)
}

type RegisterInfo struct {
	UserID int
	Conn   *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int][]*websocket.Conn),
		Register:   make(chan *RegisterInfo),
		Unregister: make(chan *RegisterInfo),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case reg := <-h.Register:
			h.mutex.Lock()
			if _, ok := h.Clients[reg.UserID]; !ok {
				h.Clients[reg.UserID] = []*websocket.Conn{}
			}
			h.Clients[reg.UserID] = append(h.Clients[reg.UserID], reg.Conn)

			// If this is the first connection, trigger Online
			if len(h.Clients[reg.UserID]) == 1 && h.OnUserOnline != nil {
				go h.OnUserOnline(reg.UserID)
			}
			h.mutex.Unlock()

		case unreg := <-h.Unregister:
			h.mutex.Lock()
			if conns, ok := h.Clients[unreg.UserID]; ok {
				for i, c := range conns {
					if c == unreg.Conn {
						h.Clients[unreg.UserID] = append(conns[:i], conns[i+1:]...)
						break
					}
				}
				if len(h.Clients[unreg.UserID]) == 0 {
					delete(h.Clients, unreg.UserID)
					// If no more connections, trigger Offline
					if h.OnUserOffline != nil {
						go h.OnUserOffline(unreg.UserID)
					}
				}
			}
			h.mutex.Unlock()

		case msg := <-h.Broadcast:
			h.mutex.Lock()
			for _, conns := range h.Clients {
				for _, conn := range conns {
					conn.WriteMessage(websocket.TextMessage, msg)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// GetOnlineUsers returns a list of UserIDs that are currently connected
func (h *Hub) GetOnlineUsers() []int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	users := make([]int, 0, len(h.Clients))
	for userID := range h.Clients {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) SendToUser(userID int, message []byte) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if conns, ok := h.Clients[userID]; ok {
		for _, conn := range conns {
			// Need to handle error?
			// if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			// 	// If write fails, we might want to unregister?
			// 	// For now simplest is just log or ignore, the loop will continue
			// }
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

var GlobalHub *Hub
