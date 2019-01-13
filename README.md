# Binar
An another demo CRUD app, with simple auth

## Requirement
- Go 1.11 (Needed for Go modules support)
- gcc (Needed to build app using SQLite3)
  - For more information about gcc, please see https://github.com/mattn/go-sqlite3
  - run `gcc -v` from terminal to check if the gcc is available in $PATH
  - gcc installation in Windows: https://github.com/mattn/go-sqlite3/issues/212#issuecomment-273531789

## Installation
- Clone this repo
  - Note: If you already define $GOPATH, please clone this repo from **outside of your $GOPATH**, because this app use Go modules and it won't work inside of $GOPATH.
- Build the app with `go build`. This command will retrive all the dependency too.
- After the app runs, it will generate `database.db` file which is an SQLite3 database file. Run the DDL file from `db/ddl_sqlite.sql` to initialize the database tables.
- Import the Postman from [here](https://www.getpostman.com/collections/433a1cdc039c288c4f5c) to play with the available endpoints.