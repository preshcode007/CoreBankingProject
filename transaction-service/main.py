#!/usr/bin/env python3
"""
Transaction Processing Service for Core Banking Application - FastAPI Version
"""

import os
import logging
from typing import Dict, List, Optional

import psycopg2
from fastapi import FastAPI, HTTPException, Query
from pydantic import BaseModel
import requests #Added this line.

from database import db  # Import the database module

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="Transaction Processing Service",
    description=(
        "Handles complex transaction processing and analysis for core banking"
    ),
    version="1.0.0"
)

# Configuration
API_HOST = os.environ.get('API_HOST', 'api')
API_PORT = os.environ.get('API_PORT', '8080')
API_BASE_URL = f"http://{API_HOST}:{API_PORT}"

# Pydantic models for request/response validation
class Transaction(BaseModel):
    id: str
    account_id: Optional[str] = None
    amount: Optional[float] = None
    type: Optional[str] = None
    status: Optional[str] = None

class TransactionResponse(BaseModel):
    transaction_id: str
    status: str
    message: str

class BatchResponse(BaseModel):
    results: List[TransactionResponse]

class AnalyticsResponse(BaseModel):
    total_transactions: int
    completed: Optional[int] = None
    pending: Optional[int] = None
    failed: Optional[int] = None
    total_deposits: Optional[float] = None
    total_withdrawals: Optional[float] = None
    net_flow: Optional[float] = None
    message: Optional[str] = None
    error: Optional[str] = None

class HealthResponse(BaseModel):
    status: str

# Database connection
conn = db.connect_db()
if not conn:
    raise Exception("Failed to connect to database")

@app.post("/process", status_code=202, response_model=Dict)
async def process_transaction(transaction: Transaction):
    """Add a transaction to the processing queue"""
    try:
        cur = conn.cursor()
        cur.execute(
            "INSERT INTO transactions (id, account_id, amount, type, status) VALUES (%s, %s, %s, %s, 'pending')",
            (transaction.id, transaction.account_id, transaction.amount, transaction.type)
        )
        conn.commit()
        cur.close()
        return {"status": "queued", "transaction_id": transaction.id}
    except psycopg2.Error as e:
        logger.error(f"Database error: {e}")
        raise HTTPException(status_code=500, detail=f"Database error: {e}")

@app.post("/process/batch", response_model=BatchResponse)
async def process_batch():
    """Process all transactions in the queue"""
    try:
        cur = conn.cursor()
        cur.execute("SELECT id, account_id, amount, type FROM transactions WHERE status = 'pending'")
        rows = cur.fetchall()
        results = []
        for row in rows:
            transaction_id, account_id, amount, transaction_type = row
            try:
                # Update status to completed
                cur.execute("UPDATE transactions SET status = 'completed' WHERE id = %s", (transaction_id,))
                conn.commit()
                results.append({
                    "transaction_id": transaction_id,
                    "status": "success",
                    "message": "Transaction processed"
                })
            except psycopg2.Error as e:
                logger.error(f"Error processing transaction {transaction_id}: {e}")
                results.append({
                    "transaction_id": transaction_id,
                    "status": "error",
                    "message": f"Database error: {e}"
                })
        cur.close()
        return {"results": results}
    except psycopg2.Error as e:
        logger.error(f"Database error: {e}")
        raise HTTPException(status_code=500, detail=f"Database error: {e}")

@app.get("/analytics", response_model=AnalyticsResponse)
async def get_analytics(account_id: Optional[str] = Query(None, description="Filter by account ID")):
    """Get transaction analytics"""
    try:
        cur = conn.cursor()
        if account_id:
            cur.execute("SELECT * FROM transactions WHERE account_id = %s", (account_id,))
        else:
            cur.execute("SELECT * FROM transactions")
        rows = cur.fetchall()
        cur.close()

        transactions = [
            {
                "id": row[0],
                "account_id": row[1],
                "amount": row[2],
                "type": row[3],
                "status": row[4]
            }
            for row in rows
        ]

        total_count = len(transactions)
        if total_count == 0:
            return {"total_transactions": 0, "message": "No transactions found"}

        completed = len([t for t in transactions if t.get('status') == 'completed'])
        pending = len([t for t in transactions if t.get('status') == 'pending'])
        failed = len([t for t in transactions if t.get('status') == 'failed'])

        deposits = sum([t.get('amount', 0) for t in transactions if t.get('type') == 'deposit' and t.get('status') == 'completed'])
        withdrawals = sum([t.get('amount', 0) for t in transactions if t.get('type') == 'withdrawal' and t.get('status') == 'completed'])

        return {
            "total_transactions": total_count,
            "completed": completed,
            "pending": pending,
            "failed": failed,
            "total_deposits": deposits,
            "total_withdrawals": withdrawals,
            "net_flow": deposits - withdrawals
        }
    except psycopg2.Error as e:
        logger.error(f"Database error: {e}")
        raise HTTPException(status_code=500, detail=f"Database error: {e}")

@app.get("/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy"}

if __name__ == '__main__':
    import uvicorn
    logger.info("Starting Transaction Processing Service...")
    uvicorn.run(app, host="0.0.0.0", port=5000)