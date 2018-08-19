package main

import (
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc"
)

import (
	pb "agent/register_client/register_service"
	"golang.org/x/net/context"
	"agent/register_client/config"
	"fmt"
	"os"
	"log"
)

func main() {
	logFile, err := os.OpenFile(config.ClientLog, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	clientLog := log.New(logFile, "[INFO]", 7)

	// TLS连接
	creds, err := credentials.NewClientTLSFromFile(config.PubKeyPath, "RegisterServer")
	if err != nil {
		clientLog.SetPrefix("[ERROR]")
		clientLog.Fatalf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", config.RemoteRpcIp, config.RemoteRpcPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		clientLog.SetPrefix("[ERROR]")
		clientLog.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewRegisterServiceClient(conn)

	// 调用方法
	reqBody := new(pb.RegisterRequest)
	reqBody.Ip = config.LocalIP
	res, err := c.NewNode(context.Background(), reqBody)
	if err != nil {
		clientLog.SetPrefix("[ERROR]")
		clientLog.Fatalln(err)
	}

	clientLog.SetPrefix("[INFO]")
	clientLog.Println("status:", res.Status, "message:", res.Message)
}
