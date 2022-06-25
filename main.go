package main

import "github.com/victorcel/crud-grpc-elasticsearch-mongodb/server"

func main() {
	quit := make(chan bool, 2)
	go server.InitializeMongoDBServer()
	go server.InitializeElasticServer()
	<-quit
}
