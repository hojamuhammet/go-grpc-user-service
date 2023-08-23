package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hojamuhammet/go-grpc-user-service/internal/utils"
	pb "github.com/hojamuhammet/go-grpc-user-service/protobuf"
)

// UserServer implements the UserServiceServer interface and provides user-related gRPC operations.
type UserServer struct {
	pb.UnimplementedUserServiceServer
	DB *sql.DB // Database connection object.
}

// GetAllUsers retrieves all users from the database and returns them as a UserList.
func (s *UserServer) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.UserList, error) {
	if req.Pagination == nil {
		req.Pagination = &pb.Pagination{
			PageSize: 10,
			PageToken: "",
		}
	}
	pageSize := int(req.Pagination.PageSize)
	pageToken := req.Pagination.PageToken

	if pageSize <= 0 {
		pageSize = 10
	}

    var query string
	var args []interface{} // Store query arguments

	if pageToken == "" {
		query = `
			SELECT id, first_name, last_name, phone_number, blocked, registration_date 
			FROM users
			LIMIT $1
		`
		args = append(args, pageSize)
	} else {
		query = `
			SELECT id, first_name, last_name, phone_number, blocked, registration_date 
			FROM users
			WHERE id > $1
			LIMIT $2
		`
		args = append(args, pageToken, pageSize)
	}

    rows, err := s.DB.Query(query, args...)
    if err != nil {
		log.Printf("Error querying users: %v", err)
        return nil, status.Error(codes.Internal, "Failed to fetch users")
    }
    defer rows.Close()

    var users []*pb.User
	var registrationDate time.Time

    for rows.Next() {
        user := &pb.User{}
        err := rows.Scan(
            &user.Id,
            &user.FirstName,
            &user.LastName,
            &user.PhoneNumber,
            &user.Blocked,
            &registrationDate,
        )
        if err != nil {
			log.Printf("Error scanning user: %v", err)
            return nil, status.Error(codes.Internal, "Failed to retrieve user data")
        }

		user.RegistrationDate = utils.ConvertToTimestamp(registrationDate)
        users = append(users, user)
    }

	var nextPageToken string
	if len(users) > 0 {
		maxID := users[len(users)-1].Id
		nextPageToken = strconv.Itoa(int(maxID))
	}

	response := &pb.UserList{
		Users: users,
		NextPageToken: nextPageToken,
	}

	log.Printf("Successfully retrieved user list")
    return response, nil
}


// GetUserById retrieves a user from the database by their ID and returns it.
func (s *UserServer) GetUserById(ctx context.Context, userID *pb.UserID) (*pb.User, error) {
	// Execute a SELECT query with a WHERE clause to fetch the user by their ID.
	user := &pb.User{}
	var registrationDate time.Time

	err := s.DB.QueryRow("SELECT * FROM users WHERE id=$1", userID.Id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Blocked,
		&registrationDate,
	)

	if err != nil {
		log.Printf("Error quiering user by ID %d: %v", userID.Id, err)
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Failed to fetch user")
	}

	user.RegistrationDate = utils.ConvertToTimestamp(registrationDate)

	log.Printf("Successfully retrieved user with ID %d", userID.Id)
	return user, nil
}

// DeleteUser deletes a user from the database by their ID and returns an empty response.
func (s *UserServer) DeleteUser(ctx context.Context, userID *pb.UserID) (*pb.Empty, error) {
	log.Printf("Deleting user with ID: %d", userID.Id)
	
	// Execute a DELETE query with a WHERE clause to remove the user with the given ID.
	result, err := s.DB.Exec("DELETE FROM users WHERE id=$1", userID.Id)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return nil, status.Error(codes.Internal, "Failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return nil, status.Error(codes.Internal, "Failed to delete user")
	}

	if rowsAffected == 0 {
		log.Printf("User not found with ID: %d", userID.Id)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	log.Printf("User with ID %d successfully deleted", userID.Id)
	return &pb.Empty{}, nil
}

func (s *UserServer) toggleBlockStatus(ctx context.Context, userID *pb.UserID, blocked bool) error {
	// Execute an UPDATE query with a WHERE clause to set the "blocked" field to the specified status for the given user ID.
	result, err := s.DB.Exec("UPDATE users SET blocked=$1 WHERE id=$2", blocked, userID.Id)
	if err != nil {
		log.Printf("Failed to update user status (UserID: %d): %v", userID.Id, err)
		return status.Error(codes.Internal, "Failed to update user status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Failed to retrieve rows affected (UserID: %d): %v", userID.Id, err)
		return status.Error(codes.Internal, "Failed to retrieve rows affected")
	}

	if rowsAffected == 0 {
		log.Printf("User not found (UserID: %d)", userID.Id)
		return status.Error(codes.NotFound, fmt.Sprintf("User with ID %d not found", userID.Id))
	}

	return nil
}

// BlockUser updates the "blocked" status of a user in the database and returns an empty response.
func (s *UserServer) BlockUser(ctx context.Context, userID *pb.UserID) (*pb.Empty, error) {
	if err := s.toggleBlockStatus(ctx, userID, true); err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("User not found: %v", err)
			return nil, status.Error(codes.NotFound, "User not found")
		}
		log.Printf("Internal server error (BlockUser): %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	log.Printf("User with ID %d successfully blocked", userID.Id)
	return &pb.Empty{}, nil
}

// UnblockUser updates the "blocked" status of a user in the database and returns an empty response.
func (s *UserServer) UnblockUser(ctx context.Context, userID *pb.UserID) (*pb.Empty, error) {
	if err := s.toggleBlockStatus(ctx, userID, false); err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		log.Printf("Internal server error (UnblockUser): %v", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}
	
	log.Printf("User with ID %d successfully unblocked", userID.Id)
	return &pb.Empty{}, nil
}

// CreateUser creates a new user in the database and returns the created user's information.
func (s *UserServer) CreateUser(ctx context.Context, userInput *pb.UserInput) (*pb.User, error) {
	var user pb.User
	var registrationDate time.Time
	
	query := `
		INSERT INTO users (first_name, last_name, phone_number)
		VALUES ($1, $2, $3)
		RETURNING id, first_name, last_name, phone_number, registration_date
	`

	err := s.DB.QueryRow(query,
		userInput.FirstName, userInput.LastName, userInput.PhoneNumber,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&registrationDate,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	user.RegistrationDate = utils.ConvertToTimestamp(registrationDate)

	log.Println("User created successfully")
	return &user, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, userUpdate *pb.UserUpdate) (*pb.User, error) {
	var user pb.User
	var registrationDate time.Time

	log.Printf("Updating user with ID: %d", userUpdate.Id)

	query := `
		UPDATE users
		SET first_name=$1, last_name=$2, phone_number=$3
		WHERE id=$4
		RETURNING *
	`

	err := s.DB.QueryRow(query,
		userUpdate.FirstName, userUpdate.LastName, userUpdate.PhoneNumber, userUpdate.Id,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Blocked,
		&registrationDate,
	)

	if err != nil {
		log.Printf("Error updating user: %v", err)
		return nil, status.Error(codes.Internal, "Failed to update user")
	}

	user.RegistrationDate = utils.ConvertToTimestamp(registrationDate)
	
	log.Println("User updated successfully")
	return &user, nil
}