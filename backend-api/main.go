package main

import (
        "encoding/json"
        "flag"
        "fmt"
        "log"
        "net/http"
        "github.com/rs/cors"
)

func main() {
        defer closeDB() // Close the database connection when the program exits.

        addr := flag.String("addr", ":8080", "HTTP network address")
        flag.Parse()

        mux := http.NewServeMux()

        mux.HandleFunc("GET /accounts", getAccounts)
        mux.HandleFunc("POST /accounts/", createAccount)
        mux.HandleFunc("GET /accounts/{id}", getAccountById)
        mux.HandleFunc("PUT /accounts/{id}", updateAccountById)
        mux.HandleFunc("DELETE /accounts/{id}", deleteAccountById)
        mux.HandleFunc("GET /transactions", getAllTransactions)
        mux.HandleFunc("POST /transactions/", createTransaction)
        mux.HandleFunc("GET /transactions/{id}", getTransactionById)
        mux.HandleFunc("PUT /transactions/{id}", updateTransactionById)

        mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
        })

        handler := cors.Default().Handler(mux)

        fmt.Printf("Banking API Server starting on port %s...\n", *addr)
        log.Fatal(http.ListenAndServe(*addr, handler))
}