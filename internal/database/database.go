package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
    "sukaina-chat/internal/models"
)

var DB *sql.DB

func Connect(databaseURL string) error {
    var err error
    DB, err = sql.Open("postgres", databaseURL)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    log.Println("✅ Database connected")
    return nil
}

func InitSchema() error {
    schema := `
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE TABLE IF NOT EXISTS users (
        user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        username VARCHAR(50) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );
    CREATE TABLE IF NOT EXISTS messages (
        message_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        sender_id UUID REFERENCES users(user_id),
        recipient_id UUID REFERENCES users(user_id),
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW()
    );`
    _, err := DB.Exec(schema)
    if err != nil {
        return fmt.Errorf("failed to create schema: %w", err)
    }
    log.Println("✅ Schema initialized")
    return nil
}

func CreateUser(username, email, passwordHash string) (*models.User, error) {
    var user models.User
    query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING user_id, username, email, created_at`
    err := DB.QueryRow(query, username, email, passwordHash).Scan(
        &user.UserID, &user.Username, &user.Email, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    query := `SELECT user_id, username, email, password_hash, created_at FROM users WHERE username = $1`
    err := DB.QueryRow(query, username).Scan(
        &user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func SaveMessage(senderID, recipientID, content string) error {
    query := `INSERT INTO messages (sender_id, recipient_id, content) VALUES ($1, $2, $3)`
    _, err := DB.Exec(query, senderID, recipientID, content)
    return err
}

func GetMessages(userID string) ([]models.Message, error) {
    query := `SELECT message_id, sender_id, recipient_id, content, created_at FROM messages WHERE recipient_id = $1 OR sender_id = $1 ORDER BY created_at DESC LIMIT 50`
    rows, err := DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var messages []models.Message
    for rows.Next() {
        var msg models.Message
        if err := rows.Scan(&msg.MessageID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.CreatedAt); err != nil {
            continue
        }
        messages = append(messages, msg)
    }
    return messages, nil
}
