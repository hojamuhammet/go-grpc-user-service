package main

import (
	"context"
	"log"

	pb "github.com/hojamuhammet/go-grpc-user-service/protobuf" // Import the generated protobuf package

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	userInput := &pb.UserInput{
		FirstName:    "Kemal",
		LastName:     "Atdayew",
		PhoneNumber:  "+993232323232",
		Password:     "K8asdasdasd!",
	}

	createdUser, err := client.CreateUser(context.Background(), userInput)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Format the timestamp for display
	registrationTime := createdUser.GetRegistrationDate()
	formattedTime := registrationTime.AsTime().Format("2006-01-02 15:04:05")

	log.Printf("User created:\n"+
		"  ID: %d\n"+
		"  First Name: %s\n"+
		"  Last Name: %s\n"+
		"  Phone Number: %s\n"+
		"  Password: %s\n"+
		"  Blocked: %t\n"+
		"  Registration Date: %s", createdUser.GetId(), createdUser.GetFirstName(), createdUser.GetLastName(),
		createdUser.GetPhoneNumber(), createdUser.GetPassword(), createdUser.GetBlocked(), formattedTime)
}
