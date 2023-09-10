package main

import (
	"fmt"
	"log"
	"net"
)

// type LogServer struct {
// 	logs.UnimplementedLogServiceServer
// 	Models data.Models
// }

// func ( l *LogServer)WriteLog(ctx, context.Context)

func (app *Config) gRPCListen() {
	_, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

}
