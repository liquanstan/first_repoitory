package main

import (
	"google.golang.org/grpc"
	"flag"
	"net"
	"fmt"
	pb "agent/command_server/command_service"
	"golang.org/x/net/context"
	"log"
	"google.golang.org/grpc/credentials"
	"os"
	"agent/command_server/cli"
	"agent/command_server/config"
)

type server struct {
	regLog *log.Logger
}

func (s *server) GetContact(ctx context.Context,in *pb.GrpcRequest)(*pb.GrpcReply, error){
	fmt.Printf("key :%v,  value :%v",in.ServerKey,in.ServerValue)
	result , err := cli.ExecuteCli("/bin/bash","-c","echo","the value is"+in.ServerValue)                             //根据接收的请求执行命令
	if err != nil {
		s.regLog.SetPrefix("[ERROR]")
		s.regLog.Println(fmt.Sprintf("excute command err :%v",err))
	} else {
		s.regLog.SetPrefix("[SUCCESS]")
		s.regLog.Println(fmt.Sprintf("excute command success! return message : %s",result))
	}
	s.regLog.SetPrefix("[SUCCESS]")
	s.regLog.Println(fmt.Sprintf("get grpc contact success"))
	out := &pb.GrpcReply{Message:"SUCCESS"}
	return out,nil
}




func main(){
	flag.Parse()
	logfile, err := os.OpenFile(config.ServerLog, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	regLog := log.New(logfile,"[INFO]",7)
	endpoint := "127.0.0.1:"+config.LocalPort
	lis,err := net.Listen("tcp",  endpoint)
	if err != nil {
		fmt.Printf("failed to listen , err:%v",err)
		regLog.SetPrefix("[ERROR]")
		regLog.Fatalln(fmt.Sprintf("failed to listen : ,%v",err))
	}

	//TLS连接
	creds,err :=credentials.NewServerTLSFromFile("command_server/tls/ca.pem","command_server/tls/server.key")
	if err != nil {
		regLog.SetPrefix("[ERROR]")
		regLog.Fatalln(fmt.Sprintf("create TLS server failed : %v",err))
	}
	grpcserver := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterContactServer(grpcserver,&server{regLog})
	err = grpcserver.Serve(lis)
	if err != nil {
		regLog.SetPrefix("[ERROR]")
		regLog.Fatalln(fmt.Sprintf("Grpc Serve failed : ,%v",err))
	}
}



