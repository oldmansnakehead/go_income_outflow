# Project Setup and Database Migration

This project uses **PostgreSQL** and **GORM** to manage the database. Before running the migrations, you need to set up an **ENUM type** in PostgreSQL to define the possible values for `TransactionType`.

### Steps to Set Up the Database:

1. **Create the `transaction_type` ENUM in PostgreSQL:**

   Before running the migration, ensure that the following **ENUM type** is created in your PostgreSQL database:

   ```sql
   CREATE TYPE transaction_type AS ENUM ('INCOME', 'EXPENSE', 'ANY');
