import axios from 'axios';

const apiBase = process.env.API_ENDPOINT;  

function axiosRequest(endpoint, method, data, authRequired = false) {
    const url = `${apiBase}${endpoint}`;
    const headers = {};
    
    if (authRequired) {
        const authToken = sessionStorage.getItem('authToken');
        if (authToken) {
            headers['Authorization'] = `Bearer ${authToken}`;
        } else {
            console.error('Error: Authentication token not found.');
            return Promise.reject(new Error('Authentication token not found'));
        }
    }
    
    return axios({ method, url, data, headers: headers }).catch(handleAxiosError);
}

function handleAxiosError(error) {
    if (error.response) {
        console.error('Error Data:', error.response.data);
        console.error('Error Status:', error.response.status);
        console.error('Error Headers:', error.response.headers);
    } else if (error.request) {
        console.error('Error Request:', error.request);
    } else {
        console.error('Error Message:', error.message);
    }
    console.error('Error Config:', error.config);
    return Promise.reject(error);
}

function submitUserInfo(username, password) {
    axiosRequest('/users', 'post', { username, password })
        .then(response => console.log('User info submitted successfully', response.data))
        .catch(error => console.error('Error submitting user info', error));
}

function authenticateUser(username, password) {
    axiosRequest('/auth', 'post', { username, password }, false)
        .then(response => {
            sessionStorage.setItem('authToken', response.data.token);
            console.log('User authenticated successfully');
        })
        .catch(error => console.error('Error authenticating user', error));
}

function submitTransactionDetail(amount, type, date) {
    axiosRequest('/transactions', 'post', { amount, type, date }, true)
        .then(response => console.log('Transaction detail submitted successfully', response.data))
        .catch(error => console.error('Error submitting transaction detail', error));
}

function submitBudgetDetail(amount, category) {
    axiosRequest('/budgets', 'post', { amount, category }, true)
        .then(response => console.log('Budget detail submitted successfully', response.data))
        .catch(error => console.error('Error submitting budget detail', error));
}

function fetchTransactions() {
    axiosRequest('/transactions', 'get', null, true)
        .then(response => displayTransactions(response.data.transactions))
        .catch(error => console.error('Error fetching transactions', error));
}

function fetchBudgets() {
    axiosRequest('/budgets', 'get', null, true)
        .then(response => displayBudgets(response.data.budgets))
        .catch(error => console.error('Error fetching budgets', error));
}

function displayTransactions(transactions) {
    transactions.forEach(transaction => console.log('Transaction:', transaction));
}

function displayBudgets(budgets) {
    budgets.forEach(budget => console.log('Budget:', budget));
}

export { submitUserInfo, authenticateUser, submitTransactionDetail, submitBudgetDetail, fetchTransactions, fetchBudgets };