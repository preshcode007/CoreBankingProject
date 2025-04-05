// API Configuration
const API_URL = process.env.API_BASE_URL || 'http://localhost:8080'; // Default for local dev
const PROCESSOR_URL = process.env.PROCESSOR_BASE_URL || 'http://localhost:5000'; // Default for local dev

// DOM Elements
const dashboardView = document.getElementById('dashboard-view');
const accountsView = document.getElementById('accounts-view');
const transactionsView = document.getElementById('transactions-view');
const dashboardLink = document.getElementById('dashboard-link');
const accountsLink = document.getElementById('accounts-link');
const transactionsLink = document.getElementById('transactions-link');

// Initialize UI
document.addEventListener('DOMContentLoaded', () => {
    // Set up navigation
    dashboardLink.addEventListener('click', showDashboard);
    accountsLink.addEventListener('click', showAccounts);
    transactionsLink.addEventListener('click', showTransactions);
    
    // Set up action buttons
    document.getElementById('new-deposit-btn').addEventListener('click', () => openTransactionModal('deposit'));
    document.getElementById('new-withdrawal-btn').addEventListener('click', () => openTransactionModal('withdrawal'));
    document.getElementById('new-account-btn').addEventListener('click', openAccountModal);
    
    // Set up form submissions
    document.getElementById('submit-transaction').addEventListener('click', submitTransaction);
    document.getElementById('submit-account').addEventListener('click', submitAccount);
    
    // Load initial data
    loadDashboardData();
});

// Navigation Functions
function showDashboard(e) {
    if (e) e.preventDefault();
    dashboardView.style.display = 'block';
    accountsView.style.display = 'none';
    transactionsView.style.display = 'none';
    dashboardLink.classList.add('active');
    accountsLink.classList.remove('active');
    transactionsLink.classList.remove('active');
    
    loadDashboardData();
}

function showAccounts(e) {
    if (e) e.preventDefault();
    dashboardView.style.display = 'none';
    accountsView.style.display = 'block';
    transactionsView.style.display = 'none';
    dashboardLink.classList.remove('active');
    accountsLink.classList.add('active');
    transactionsLink.classList.remove('active');
    
    loadAccounts();
}

function showTransactions(e) {
    if (e) e.preventDefault();
    dashboardView.style.display = 'none';
    accountsView.style.display = 'none';
    transactionsView.style.display = 'block';
    dashboardLink.classList.remove('active');
    accountsLink.classList.remove('active');
    transactionsLink.classList.add('active');
    
    loadTransactions();
}

// Data Loading Functions
async function loadDashboardData() {
    try {
        // Load accounts for balance calculation
        const accountsResponse = await fetch(`${API_URL}/accounts`);
        const accounts = await accountsResponse.json();
        
        // Calculate total balance
        const totalBalance = accounts.reduce((total, account) => total + account.balance, 0);
        document.getElementById('total-balance').textContent = formatCurrency(totalBalance);
        
        // Load recent transactions
        const transactionsResponse = await fetch(`${API_URL}/transactions`);
        const transactions = await transactionsResponse.json();
        
        // Sort by most recent (assuming ID is sequential)
        const recentTransactions = transactions
            .sort((a, b) => parseInt(b.id) - parseInt(a.id))
            .slice(0, 5);
            
        displayRecentTransactions(recentTransactions);
        
        // Get analytics
        const analyticsResponse = await fetch(`${PROCESSOR_URL}/analytics`);
        const analytics = await analyticsResponse.json();
        console.log('Analytics:', analytics);
        
    } catch (error) {
        console.error('Error loading dashboard data:', error);
        alert('Failed to load dashboard data. Please check the console for details.');
    }
}

async function loadAccounts() {
    try {
        const response = await fetch(`${API_URL}/accounts`);
        const accounts = await response.json();
        
        const accountsList = document.getElementById('accounts-list');
        if (accounts.length === 0) {
            accountsList.innerHTML = '<p class="text-muted">No accounts found</p>';
            return;
        }
        
        accountsList.innerHTML = accounts.map(account => `
            <div class="list-group-item">
                <div class="d-flex w-100 justify-content-between">
                    <h5 class="mb-1">Account #${account.id}</h5>
                    <span class="fw-bold">${formatCurrency(account.balance)}</span>
                </div>
                <p class="mb-1">Owner: ${account.owner}</p>
            </div>
        `).join('');
        
    } catch (error) {
        console.error('Error loading accounts:', error);
        alert('Failed to load accounts. Please check the console for details.');
    }
}

async function loadTransactions() {
    try {
        const response = await fetch(`${API_URL}/transactions`);
        const transactions = await response.json();
        
        const transactionsList = document.getElementById('transactions-list');
        if (transactions.length === 0) {
            transactionsList.innerHTML = '<p class="text-muted">No transactions found</p>';
            return;
        }
        
        // Sort by most recent (assuming ID is sequential)
        const sortedTransactions = transactions.sort((a, b) => parseInt(b.id) - parseInt(a.id));
        
        transactionsList.innerHTML = sortedTransactions.map(txn => `
            <div class="list-group-item">
                <div class="d-flex w-100 justify-content-between">
                    <h5 class="mb-1">${txn.type.charAt(0).toUpperCase() + txn.type.slice(1)}</h5>
                    <span class="fw-bold ${txn.type === 'deposit' ? 'positive' : 'negative'}">
                        ${txn.type === 'deposit' ? '+' : '-'}${formatCurrency(txn.amount)}
                    </span>
                </div>
                <p class="mb-1">Account: #${txn.account_id}</p>
                <p class="mb-1">${txn.description}</p>
                <small class="text-muted">Status: ${txn.status}</small>
            </div>
        `).join('');
        
    } catch (error) {
        console.error('Error loading transactions:', error);
        alert('Failed to load transactions. Please check the console for details.');
    }
}

function displayRecentTransactions(transactions) {
    const recentTransactionsEl = document.getElementById('recent-transactions');
    
    if (transactions.length === 0) {
        recentTransactionsEl.innerHTML = '<p class="text-muted">No recent transactions</p>';
        return;
    }
    
    recentTransactionsEl.innerHTML = transactions.map(txn => `
        <div class="mb-3">
            <div class="d-flex justify-content-between">
                <span>${txn.description}</span>
                <span class="${txn.type === 'deposit' ? 'positive' : 'negative'}">
                    ${txn.type === 'deposit' ? '+' : '-'}${formatCurrency(txn.amount)}
                </span>
            </div>
            <small class="text-muted">Account #${txn.account_id} • Status: ${txn.status}</small>
        </div>
    `).join('<hr>');
}

// Modal Functions
async function openTransactionModal(type) {
    // Set transaction type
    document.getElementById('transaction-type').value = type;
    document.getElementById('transaction-modal-title').textContent = 
        `New ${type.charAt(0).toUpperCase() + type.slice(1)}`;
// Load accounts for dropdown
    try {
        const response = await fetch(`${API_URL}/accounts`);
        const accounts = await response.json();
        
        const accountSelect = document.getElementById('account-select');
        accountSelect.innerHTML = '<option value="">Select an account</option>';
        
        accounts.forEach(account => {
            const option = document.createElement('option');
            option.value = account.id;
            option.textContent = `#${account.id} - ${account.owner} (${formatCurrency(account.balance)})`;
            accountSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading accounts:', error);
    }
    
    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('transaction-modal'));
    modal.show();
}

function openAccountModal() {
    // Reset form
    document.getElementById('account-form').reset();
    
    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('account-modal'));
    modal.show();
}

// Form Submission Functions
async function submitTransaction() {
    const accountId = document.getElementById('account-select').value;
    const amount = document.getElementById('amount-input').value;
    const description = document.getElementById('description-input').value;
    const type = document.getElementById('transaction-type').value;
    
    if (!accountId || !amount || !description) {
        alert('Please fill out all fields');
        return;
    }
    
    try {
        // Create transaction
        const createResponse = await fetch(`${API_URL}/transactions`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                account_id: accountId,
                amount: parseFloat(amount),
                type: type,
                description: description
            })
        });
        
        if (!createResponse.ok) {
            throw new Error(`Failed to create transaction: ${createResponse.statusText}`);
        }
        
        const transaction = await createResponse.json();
        
        // Process transaction
        await fetch(`${PROCESSOR_URL}/process`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(transaction)
        });
        
        // Close modal
        bootstrap.Modal.getInstance(document.getElementById('transaction-modal')).hide();
        
        // Reload dashboard data
        loadDashboardData();
        
        alert(`${type.charAt(0).toUpperCase() + type.slice(1)} processed successfully`);
    } catch (error) {
        console.error('Error creating transaction:', error);
        alert('Failed to process transaction. Please check the console for details.');
    }
}

async function submitAccount() {
    const owner = document.getElementById('owner-input').value;
    const initialBalance = document.getElementById('initial-balance-input').value;
    
    if (!owner || !initialBalance) {
        alert('Please fill out all fields');
        return;
    }
    
    try {
        const response = await fetch(`${API_URL}/accounts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                owner: owner,
                balance: parseFloat(initialBalance)
            })
        });
        
        if (!response.ok) {
            throw new Error(`Failed to create account: ${response.statusText}`);
        }
        
        // Close modal
        bootstrap.Modal.getInstance(document.getElementById('account-modal')).hide();
        
        // Reload dashboard data
        loadDashboardData();
        
        alert('Account created successfully');
    } catch (error) {
        console.error('Error creating account:', error);
        alert('Failed to create account. Please check the console for details.');
    }
}

// Utility Functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount);
}
