package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
    dbHost := os.Getenv("DATABASE_HOST")
    dbUser := os.Getenv("DATABASE_USER")
    dbPassword := os.Getenv("DATABASE_PASSWORD")
    dbName := os.Getenv("DATABASE_NAME")
    dbPort := os.Getenv("DATABASE_PORT")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }

    fmt.Println("Successfully connected to PostgreSQL!")

    return db
}