import axios from 'axios';

function submitUserInfo(username, password) {
    const apiEndpoint = process.env.API_ENDPOINT + '/users';
    axios.post(apiEndpoint, {
        username: username,
        password: password
    }).then(response => {
        console.log('User info submitted successfully', response.data);
    }).catch(error => {
        console.error('Error submitting user info', error);
    });
}

function authenticateUser(username, password) {
    const apiEndpoint = process.env.API_ENDPOINT + '/auth';
    axios.post(apiEndpoint, {
        username: username,
        password: password
    }).then(response => {
        sessionStorage.setItem('authToken', response.data.token);
        console.log('User authenticated successfully');
    }).catch(error => {
        console.error('Error authenticating user', error);
    });
}

function submitTransactionDetail(amount, type, date) {
    const apiEndpoint = process.env.API_ENDPOINT + '/transactions';
    const authToken = sessionStorage.getItem('authToken');
    axios.post(apiEndpoint, {
        amount: amount,
        type: type,
        date: date
    }, {
        headers: {
            'Authorization': `Bearer ${authToken}`
        }
    }).then(response => {
        console.log('Transaction detail submitted successfully', response.data);
    }).catch(error => {
        console.error('Error submitting transaction detail', error);
    });
}

function submitBudgetDetail(amount, category) {
    const apiEndpoint = process.env.API_ENDPOINT + '/budgets';
    const authToken = sessionStorage.getItem('authToken');
    axios.post(apiEndpoint, {
        amount: amount,
        category: category
    }, {
        headers: {
            'Authorization': `Bearer ${authToken}`
        }
    }).then(response => {
        console.log('Budget detail submitted successfully', response.data);
    }).catch(error => {
        console.error('Error submitting budget detail', error);
    });
}

function fetchTransactions() {
    const apiEndpoint = process.env.API_ENDPOINT + '/transactions';
    const authToken = sessionStorage.getItem('authToken');
    axios.get(apiEndpoint, {
        headers: {
            'Authorization': `Bearer ${authToken}`
        }
    }).then(response => {
        displayTransactions(response.data.transactions);
    }).catch(error => {
        console.error('Error fetching transactions', error);
    });
}

function displayTransactions(transactions) {
    transactions.forEach(transaction => {
        console.log('Transaction:', transaction);
    });
}

function fetchBudgets() {
    const apiEndpoint = process.env.API_ENDPOINT + '/budgets';
    const authToken = sessionStorage.getItem('authToken');
    axios.get(apiEndpoint, {
        headers: {
            'Authorization': `Bearer ${authToken}`
        }
    }).then(response => {
        displayBudgets(response.data.budgets);
    }).catch(error => {
        console.error('Error fetching budgets', error);
    });
}

function displayBudgets(budgets) {
    budgets.forEach(budget => {
        console.log('Budget:', budget);
    });
}

export { submitUserInfo, authenticateUser, submitTransactionDetail, submitBudgetDetail, fetchTransactions, fetchBudgets };
