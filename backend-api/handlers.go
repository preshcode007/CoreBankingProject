package main

import (
        "backend-api/database"
        "database/sql"
        "encoding/json"
        "net/http"
        "strconv"
)

// Global database variable
var db *sql.DB

// Initialize Database Connection. This should be done once when the program starts.
func init() {
        db = database.ConnectDB()
}

// Ensure the database connection is closed when the program exits.
func closeDB() {
        if db != nil {
                db.Close()
        }
}

func getAccounts(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        rows, err := db.Query("SELECT id, balance FROM accounts")
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        defer rows.Close()

        var accounts []Account
        for rows.Next() {
                var account Account
                if err := rows.Scan(&account.ID, &account.Balance); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                accounts = append(accounts, account)
        }

        json.NewEncoder(w).Encode(accounts)
}

func createAccount(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var newAccount Account
        err := json.NewDecoder(r.Body).Decode(&newAccount)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        var id int
        err = db.QueryRow("INSERT INTO accounts (balance) VALUES ($1) RETURNING id", newAccount.Balance).Scan(&id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        newAccount.ID = strconv.Itoa(id)
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(newAccount)
}

func getAccountById(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.PathValue("id")

        var account Account
        err := db.QueryRow("SELECT id, balance FROM accounts WHERE id = $1", id).Scan(&account.ID, &account.Balance)
        if err != nil {
                if err == sql.ErrNoRows {
                        http.Error(w, "Account not found", http.StatusNotFound)
                } else {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                return
        }

        json.NewEncoder(w).Encode(account)
}

func updateAccountById(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.PathValue("id")

        var updatedAccount Account
        err := json.NewDecoder(r.Body).Decode(&updatedAccount)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        _, err = db.Exec("UPDATE accounts SET balance = $1 WHERE id = $2", updatedAccount.Balance, id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        updatedAccount.ID = id
        json.NewEncoder(w).Encode(updatedAccount)
}

func deleteAccountById(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.PathValue("id")

        _, err := db.Exec("DELETE FROM accounts WHERE id = $1", id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        w.WriteHeader(http.StatusNoContent)
}

func getAllTransactions(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        rows, err := db.Query("SELECT id, account_id, amount, type, status FROM transactions")
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        defer rows.Close()

        var transactions []Transaction
        for rows.Next() {
                var transaction Transaction
                if err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Type, &transaction.Status); err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                transactions = append(transactions, transaction)
        }

        json.NewEncoder(w).Encode(transactions)
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var newTxn Transaction
        err := json.NewDecoder(r.Body).Decode(&newTxn)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        // Verify account exists
        var accountExists bool
        err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", newTxn.AccountID).Scan(&accountExists)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        if !accountExists {
                http.Error(w, "Account not found", http.StatusBadRequest)
                return
        }

        // Process transaction
        var id int
        err = db.QueryRow("INSERT INTO transactions (account_id, amount, type, status) VALUES ($1, $2, $3, 'pending') RETURNING id", newTxn.AccountID, newTxn.Amount, newTxn.Type).Scan(&id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
        newTxn.ID = strconv.Itoa(id)

        // Update account balance
        if newTxn.Type == "deposit" {
                _, err = db.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", newTxn.Amount, newTxn.AccountID)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                _, err = db.Exec("UPDATE transactions SET status = 'completed' WHERE id = $1", newTxn.ID)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
        } else if newTxn.Type == "withdrawal" {
                var balance float64
                err = db.QueryRow("SELECT balance FROM accounts WHERE id = $1", newTxn.AccountID).Scan(&balance)
                if err != nil {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                        return
                }
                if balance >= newTxn.Amount {
                        _, err = db.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2", newTxn.Amount, newTxn.AccountID)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        _, err = db.Exec("UPDATE transactions SET status = 'completed' WHERE id = $1", newTxn.ID)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                } else {
                        _, err = db.Exec("UPDATE transactions SET status = 'failed' WHERE id = $1", newTxn.ID)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                }
        }

        json.NewEncoder(w).Encode(newTxn)
}

func getTransactionById(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.PathValue("id")

        var transaction Transaction
        err := db.QueryRow("SELECT id, account_id, amount, type, status FROM transactions WHERE id = $1", id).Scan(&transaction.ID, &transaction.AccountID, &transaction.Amount, &transaction.Type, &transaction.Status)
        if err != nil {
                if err == sql.ErrNoRows {
                        http.Error(w, "Transaction not found", http.StatusNotFound)
                } else {
                        http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                return
        }

        json.NewEncoder(w).Encode(transaction)
}

func updateTransactionById(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.PathValue("id")

        var updatedTxn Transaction
        err := json.NewDecoder(r.Body).Decode(&updatedTxn)
        if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
        }

        _, err = db.Exec("UPDATE transactions SET status = $1 WHERE id = $2", updatedTxn.Status, id)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        json.NewEncoder(w).Encode(updatedTxn)
}