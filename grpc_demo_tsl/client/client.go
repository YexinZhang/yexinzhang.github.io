package main

import (
	"context"
	"fmt"
	"go_dev/grpc_demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
)

func main() {

	creds,err :=  credentials.NewClientTLSFromFile("/home/zhangyexin/ssl/grpc.crt","grpc.com")
	fmt.Println(creds)
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := grpc.Dial("127.0.0.1:7891", grpc.WithTransportCredentials(creds))
	if err != nil {
		fmt.Println("dial: err:",err)
		return
	}
	sendClient := pb.NewSendTestMsgClient(conn)
	stream, err := sendClient.SendT(context.Background(), &pb.TestReQ{Label: "grpc"})
	if err != nil {
		fmt.Println("send err: ",err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Println(err)
			if err == io.EOF {
				break
			}
			continue
		}
		fmt.Println(resp.Msg, resp.Status)
	}
}
