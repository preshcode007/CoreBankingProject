import pytest
import requests
import os

API_HOST = os.environ.get('API_HOST', '127.0.0.1')  # Use localhost for testing
API_PORT = os.environ.get('API_PORT', '5000')
API_BASE_URL = f"http://{API_HOST}:{API_PORT}"

def test_health_check():
    response = requests.get(f"{API_BASE_URL}/health")
    assert response.status_code == 200
    assert response.json() == {"status": "healthy"}

def test_process_transaction():
    transaction_data = {
        "id": "test_txn_1",
        "account_id": "test_account_1",
        "amount": 100.0,
        "type": "deposit",
    }
    response = requests.post(f"{API_BASE_URL}/process", json=transaction_data)
    assert response.status_code == 202
    assert response.json()["status"] == "queued"

def test_process_batch():
    # Assuming you have some pending transactions in the database from previous tests
    response = requests.post(f"{API_BASE_URL}/process/batch")
    assert response.status_code == 200
    assert "results" in response.json()
    #add more assertions here based on your database state.

def test_analytics():
    response = requests.get(f"{API_BASE_URL}/analytics")
    assert response.status_code == 200
    assert "total_transactions" in response.json()

def test_analytics_with_account_id():
    response = requests.get(f"{API_BASE_URL}/analytics?account_id=test_account_1") #add an account id that exists in your database.
    assert response.status_code == 200
    assert "total_transactions" in response.json()

# Add more tests as needed