# go_income_outflow

A Go project to manage income and expenditure data with integration to a PostgreSQL database. This project includes functionality for tracking transactions, user accounts, credit card debts, and more.

## Features

- **Transaction Tracking**: Track different types of transactions such as income, expense, and more.
- **User Management**: Handle user accounts, authentication, and relationships to transactions.
- **Credit Card Debt Management**: Track credit card debts and their statuses.
- **Database Integration**: Integration with PostgreSQL for managing data.
- **Flexible ENUM Management**: Supports dynamic creation of ENUM types in PostgreSQL to categorize transaction types and other entities.
- **Migrations**: Auto-migration of database schemas for easy setup and updates.

## Setup

### Prerequisites

- Go (v1.16 or later)
- PostgreSQL (v12 or later)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/oldmansnakehead/go_income_outflow.git
   cd go_income_outflow

2. Install dependencies:
   ```bash
   go mod tidy
3. Set up your .env file (for database credentials):
   - Create a .env file in the root of the project and provide the following environment variables
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=go_income_outflow
   DB_SSL_MODE=disable
4. Run the application:
   - Start the application
   ```bash
   go run main.go
5. The server will start, and you can access the API.

## API Endpoints
### /users
   - GET /users: Get all users.
   - POST /users: Create a new user.
### /transactions
   - GET /transactions: Get all transactions.
   - POST /transactions: Create a new transaction.

### /accounts
   - GET /accounts: Get all accounts.
   - POST /accounts: Create a new account.

### /creditcarddebt
   - GET /creditcarddebt: Get all credit card debts.
   - POST /creditcarddebt: Create a new credit card debt.

## Database Schema

This project uses PostgreSQL to store data. The following tables and relationships exist:

   - **Users**: Stores user information.
   - **Accounts**: Stores user accounts.
   - **Transactions**: Stores transactions related to income and expenses.
   - **CreditCardDebt**: Stores information about credit card debts.
   - **TransactionCategory**: Stores categories for transactions (e.g., groceries, entertainment).
