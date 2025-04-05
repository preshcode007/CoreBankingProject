package main

// Account represents a bank account from the database
type Account struct {
        ID      string  `json:"id"`
        Balance float64 `json:"balance"`
}

// Transaction represents a banking transaction from the database
type Transaction struct {
        ID          string  `json:"id"`
        AccountID   string  `json:"account_id"`
        Amount      float64 `json:"amount"`
        Type        string  `json:"type"`        // "deposit", "withdrawal"
        Description string  `json:"description"` //If you want to add description.
        Status      string  `json:"status"`      // "pending", "completed", "failed"
}