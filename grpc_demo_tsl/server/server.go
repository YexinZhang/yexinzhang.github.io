package main

import (
	"fmt"
	"go_dev/grpc_demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"net"
)

type msg struct {}
func (m *msg) SendT(q *pb.TestReQ, spres pb.SendTestMsg_SendTServer) error {
	label := q.Label
	fmt.Println(label)
	for i := 1; i <= 100; i++ {
		err := spres.Send(&pb.ResponsE{
			Msg:    fmt.Sprintf("msg: %d\n", i),
			Status: fmt.Sprintf("status: %d",i),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	listen, err := net.Listen("tcp", "0.0.0.0:7891")
	if err != nil {
		fmt.Println(err)
		return
	}
	tls, err := credentials.NewServerTLSFromFile("/home/zhangyexin/ssl/grpc.crt","/home/zhangyexin/ssl/grpc.key")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := grpc.NewServer(grpc.Creds(tls))
	pb.RegisterSendTestMsgServer(server, &msg{})
	reflection.Register(server)
	err = server.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}
