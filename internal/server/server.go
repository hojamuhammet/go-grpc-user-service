package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	pb "github.com/hojamuhammet/go-grpc-user-service/protobuf"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserServer implements the UserServiceServer interface and provides user-related gRPC operations.
type UserServer struct {
	pb.UnimplementedUserServiceServer
	DB *sql.DB // Database connection object.
}

// GetAllUsers retrieves all users from the database and returns them as a UserList.
func (s *UserServer) GetAllUsers(ctx context.Context, empty *pb.Empty) (*pb.UserList, error) {
	query := "SELECT id, first_name, last_name, phone_number, password, blocked, registration_date FROM users"
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User

	for rows.Next() {
		user := &pb.User{}
		var registrationTime time.Time

		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.PhoneNumber,
			&user.Password,
			&user.Blocked,
			&registrationTime,
		)
		if err != nil {
			return nil, err
		}

		user.RegistrationDate = timestamppb.New(registrationTime)
		user.RegistrationDateString = registrationTime.Format("02-01-2006 15:04:05 MST")
		users = append(users, user)
	}

	return &pb.UserList{Users: users}, nil
}



// GetUserById retrieves a user from the database by their ID and returns it.
func (s *UserServer) GetUserById(ctx context.Context, userID *pb.UserID) (*pb.User, error) {
	// Execute a SELECT query with a WHERE clause to fetch the user by their ID.
	user := &pb.User{}

	var registrationTime time.Time

	err := s.DB.QueryRow("SELECT * FROM users WHERE id=$1", userID.Id).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Password,
		&user.Blocked,
		&registrationTime,
	)

	if err != nil {
		return nil, err
	}

	// Format the registrationTime as "yyyy:mm:dd hh:mm:ss".
	registrationTimeString := registrationTime.Format("2006-01-02 15:04:05")

	// Set the formatted registration date in the registration_date_string field.
	user.RegistrationDateString = registrationTimeString

	return user, nil
}

// DeleteUser deletes a user from the database by their ID and returns an empty response.
func (s *UserServer) DeleteUser(ctx context.Context, userID *pb.UserID) (*pb.Empty, error) {
	// Execute a DELETE query with a WHERE clause to remove the user with the given ID.
	_, err := s.DB.Exec("DELETE FROM users WHERE id=$1", userID.Id)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

// BlockUser updates the "blocked" status of a user in the database and returns an empty response.
func (s *UserServer) BlockUser(ctx context.Context, userID *pb.UserID) (*pb.Empty, error) {
	// Execute an UPDATE query with a WHERE clause to set the "blocked" field to true for the given user ID.
	result, err := s.DB.Exec("UPDATE users SET blocked=true WHERE id=$1", userID.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("user with ID %d not found", userID.Id)
	}

	return &pb.Empty{}, nil
}

// CreateUser creates a new user in the database and returns the created user's information.
func (s *UserServer) CreateUser(ctx context.Context, userInput *pb.UserInput) (*pb.User, error) {
	var user pb.User

	registrationTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	query := `
		INSERT INTO users (first_name, last_name, phone_number, password, registration_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *
	`

	err := s.DB.QueryRow(query,
		userInput.FirstName, userInput.LastName, userInput.PhoneNumber, userInput.Password, registrationTime,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Password,
		&user.Blocked,
		&registrationTime,
	)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	registrationTimeParsed, _ := time.Parse("2006-01-02 15:04:05", registrationTime)
	user.RegistrationDate = timestamppb.New(registrationTimeParsed)
	log.Println("User created successfully")
	return &user, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, userUpdate *pb.UserUpdate) (*pb.User, error) {
	var (
		user pb.User
		registrationTime time.Time
	)

	// Generate the new registration time.
	registrationTime = time.Now()

	query := `
		UPDATE users
		SET first_name=$1, last_name=$2, phone_number=$3, password=$4, blocked=$5, registration_date=$6
		WHERE id=$7
		RETURNING *
	`

	err := s.DB.QueryRow(query,
		userUpdate.FirstName, userUpdate.LastName, userUpdate.PhoneNumber, userUpdate.Password, userUpdate.Blocked, registrationTime.Format("2006-01-02 15:04:05"), userUpdate.Id,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Password,
		&user.Blocked,
		&registrationTime,
	)

	if err != nil {
		return nil, err
	}

	user.RegistrationDate = timestamppb.New(registrationTime)
	user.RegistrationDateString = registrationTime.Format("02-01-2006 15:04:05 MST")
	return &user, nil
}
