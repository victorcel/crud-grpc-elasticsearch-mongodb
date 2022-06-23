package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/database"
	userpb "github.com/victorcel/crud-grpc-elasticsearch-mongodb/proto"
	"github.com/victorcel/crud-grpc-elasticsearch-mongodb/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Iniciado servidor...")

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listenServer, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	repository, err := database.NewMongoDbRepository(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}

	serverMain := server.NewServer(repository)

	serverGrpc := grpc.NewServer()
	userpb.RegisterUserServiceServer(serverGrpc, serverMain)

	reflection.Register(serverGrpc)

	fmt.Println("Server run 8000")

	if err := serverGrpc.Serve(listenServer); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}

}
