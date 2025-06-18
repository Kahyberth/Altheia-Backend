package websocket

import (
	"encoding/json"
	"log"
	"time"
)

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("WebSocket client registered: %s for clinic: %s", client.ID, client.ClinicID)

			welcomeMsg := Message{
				Type:      "connection",
				Data:      map[string]string{"status": "connected", "client_id": client.ID},
				Timestamp: time.Now(),
				ClinicID:  client.ClinicID,
			}
			if data, err := json.Marshal(welcomeMsg); err == nil {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}

		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("WebSocket client unregistered: %s", client.ID)
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) BroadcastToClinic(clinicID string, message []byte) {
	for client := range h.Clients {
		if client.ClinicID == clinicID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
}

func (h *Hub) GetConnectedClients() int {
	return len(h.Clients)
}

func (h *Hub) GetClinicClients(clinicID string) int {
	count := 0
	for client := range h.Clients {
		if client.ClinicID == clinicID {
			count++
		}
	}
	return count
}
