package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
    "sukaina-chat/internal/auth"
    "sukaina-chat/internal/database"
    "sukaina-chat/internal/models"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins (change in production)
    },
}

type Client struct {
    UserID string
    Conn   *websocket.Conn
    Send   chan []byte
}

type Hub struct {
    clients    map[string]*Client
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client.UserID] = client
            h.mutex.Unlock()
            log.Printf("✅ Client connected: %s", client.UserID)

        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client.UserID]; ok {
                delete(h.clients, client.UserID)
                close(client.Send)
            }
            h.mutex.Unlock()
            log.Printf("❌ Client disconnected: %s", client.UserID)
        }
    }
}

func (h *Hub) SendToUser(userID string, message []byte) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    if client, ok := h.clients[userID]; ok {
        select {
        case client.Send <- message:
        default:
            close(client.Send)
            delete(h.clients, userID)
        }
    }
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    tokenString := r.URL.Query().Get("token")
    claims, err := auth.ValidateToken(tokenString)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }

    client := &Client{
        UserID: claims.UserID,
        Conn:   conn,
        Send:   make(chan []byte, 256),
    }

    h.register <- client

    go client.writePump()
    go client.readPump(h)
}

func (c *Client) readPump(h *Hub) {
    defer func() {
        h.unregister <- c
        c.Conn.Close()
    }()

    for {
        _, message, err := c.Conn.ReadMessage()
        if err != nil {
            break
        }

        var wsMsg models.WebSocketMessage
        if err := json.Unmarshal(message, &wsMsg); err != nil {
            continue
        }

        wsMsg.SenderID = c.UserID
        wsMsg.Timestamp = time.Now().Unix()

        // Save to database
        database.SaveMessage(c.UserID, wsMsg.RecipientID, wsMsg.Content)

        // Send to recipient
        responseBytes, _ := json.Marshal(wsMsg)
        h.SendToUser(wsMsg.RecipientID, responseBytes)
    }
}

func (c *Client) writePump() {
    defer c.Conn.Close()

    for message := range c.Send {
        if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
            break
        }
    }
}
