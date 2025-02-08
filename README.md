# Everlock

## Description

Everlock is a secure message encryption platform that ensures a secret can only be revealed when all designated participants enter their unique key fragments. Using distributed encryption, Everlock enables users to create and store encrypted "last wishes" or confidential messages, which can only be unlocked through collective decryption.

## 🔒 Key Features:

*   Multi-key encryption for ultimate security
*   Decentralized access control—no single person can decrypt the message alone
*   Simple and intuitive interface powered by Mustache
*   Built with Go for performance and reliability
*   Ideal for safeguarding final messages, secret agreements, or time-sensitive information

## Linter

To check the code for style and standard compliance, you can use the linter via the `golangci-lint` tool.

### Installing the Linter

Install `golangci-lint` if it's not already installed:
```bash
brew install golangci-lint
```

### Running the Linter

To run the linter, execute the following command in the root of the project:
```bash
golangci-lint run
```
The linter will check all Go files in the project for errors and style warnings.

## Migrations

To apply migrations to the database, use the golang-migrate tool.

### Installing migration tool

Ensure you have golang-migrate installed. To install it, run the following command:
```bash
brew install golang-migrate
```

### Running migration tool 

Apply the migrations with this command:
```bash
migrate -path internal/db/migrations -database "postgres://everlock:12345@localhost:5432/everlock?sslmode=disable" up
```
Where:
- path internal/db/migrations is the path to the migrations folder.
- database "postgres://everlock:12345@localhost:5432/everlock?sslmode=disable" is the connection string for your database.

To roll back migrations, use this command:
```bash
migrate -path internal/db/migrations -database "postgres://everlock:12345@localhost:5432/everlock?sslmode=disable" down
```

To check the current status of migrations, run:
```bash
migrate -path internal/db/migrations -database "postgres://everlock:12345@localhost:5432/everlock?sslmode=disable" status
```
