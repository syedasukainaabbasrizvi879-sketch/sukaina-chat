```go
package models

import "time"

type User struct {
    UserID       string    `json:"user_id"`
    Username     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"-"`
    CreatedAt    time.Time `json:"created_at"`
}

type Message struct {
    MessageID   string    `json:"message_id"`
    SenderID    string    `json:"sender_id"`
    RecipientID string    `json:"recipient_id"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"created_at"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegisterRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token    string `json:"token"`
    UserID   string `json:"user_id"`
    Username string `json:"username"`
}

type WebSocketMessage struct {
    Type        string `json:"type"`
    RecipientID string `json:"recipient_id"`
    Content     string `json:"content"`
    SenderID    string `json:"sender_id,omitempty"`
    Timestamp   int64  `json:"timestamp,omitempty"`
}
```

---
