# *gRPC CRUD Application*

This repository conatains a gRPC-based CRUD application built in Go. 
The application provides basic Create, Read, Update, and Delete (CRUD) operations for managing user data. 
It uses PostgreSQL as the database and gRPC for communication.

## *Prerequisites<br/>*
Before running the application, make sure you have the following prerequisites:<br/>
- Go (1.16 or higher)
- PostgreSQL database
- Protocol Buffers compiler (protoc) for generating Go code from .proto files
Installation and Setup
Clone this repository to your local machine:</br>
```cmd
git clone https://github.com/hojamuhammet/grpc-crud-go.git
```
```cmd
cd grpc-crud-go 
```
## Install the required Go packages and dependencies:

- Set up your PostgreSQL database and create the users table. You can use the provided user.sql file to create the required table structure.
- Create a .env file in the project root and set the following environment variables:
  - DB_HOST=<your_db_host>
  - DB_PORT=<your_db_port>
  - DB_USER=<your_db_user>
  - DB_PASSWORD=<your_db_password>
  - DB_NAME=<your_db_name>
  - GRPC_PORT=<desired_grpc_port>

## *Compile the .proto file to generate Go code:</br>*

```cmd
protoc -I=protobuf --go_out=. --go-grpc_out=. user_proto.proto
```

## *Build and run the application:*

```cmd
go run main.go
```

## API Documentation
The application provides the following gRPC service methods:
- GetAllUsers: Retrieves a list of all users.
- GetUserById: Retrieves a single user by their ID.
- CreateUser: Creates a new user with the provided input data.
- UpdateUser: Updates an existing user's data
- DeleteUser: Deletes a user by their ID.
- BlockUser: Blocks a user by their ID.
- Refer to the user_proto.proto file for detailed method and message definitions.
