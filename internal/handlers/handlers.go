```go
package handlers

import (
    "encoding/json"
    "net/http"

    "sukaina-chat/internal/auth"
    "sukaina-chat/internal/database"
    "sukaina-chat/internal/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var req models.RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }

    if req.Username == "" || req.Email == "" || req.Password == "" {
        http.Error(w, `{"error":"all fields required"}`, http.StatusBadRequest)
        return
    }

    hashedPassword, err := auth.HashPassword(req.Password)
    if err != nil {
        http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
        return
    }

    user, err := database.CreateUser(req.Username, req.Email, hashedPassword)
    if err != nil {
        http.Error(w, `{"error":"user already exists"}`, http.StatusConflict)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "registration successful",
        "user_id": user.UserID,
    })
}

func Login(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
        return
    }

    user, err := database.GetUserByUsername(req.Username)
    if err != nil {
        http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    if !auth.CheckPassword(req.Password, user.PasswordHash) {
        http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateToken(user.UserID, user.Username)
    if err != nil {
        http.Error(w, `{"error":"failed to generate token"}`, http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(models.LoginResponse{
        Token:    token,
        UserID:   user.UserID,
        Username: user.Username,
    })
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    tokenString := r.Header.Get("Authorization")
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }

    claims, err := auth.ValidateToken(tokenString)
    if err != nil {
        http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
        return
    }

    messages, err := database.GetMessages(claims.UserID)
    if err != nil {
        http.Error(w, `{"error":"failed to get messages"}`, http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "messages": messages,
    })
}

func Health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"healthy"}`))
}
```

---

# PART 4: WEBSOCKET (Real-time Chat)

---
