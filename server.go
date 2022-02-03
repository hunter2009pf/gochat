package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	pb "gochat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	port = flag.Int("port", 8888, "The server port")
)

type server struct {
	pb.UnimplementedAuthServiceServer
}

type commServer struct {
	pb.UnimplementedCommunicationServer
	connList []*Connection
}

type Connection struct {
	stream pb.Communication_ConnectServerServer
	userId string
	err    chan error
}

func (s *server) GetIMToken(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Printf("Received: %v", in.GetUserId())
	return &pb.AuthResponse{ImToken: in.GetUserId() + strconv.Itoa(int(time.Now().UnixMilli()))}, nil
}

func (s *commServer) ConnectServer(user *pb.User, stream pb.Communication_ConnectServerServer) error {
	conn := &Connection{
		stream: stream,
		userId: user.GetUserId(),
		err:    make(chan error),
	}
	s.connList = append(s.connList, conn)
	return <-conn.err
}

func (s *commServer) SendMsg(ctx context.Context, in *pb.SendMsgRequest) (*pb.SendMsgResponse, error) {
	log.Print("start sending msg")
	if s.connList == nil {
		log.Print("connList is nil")
		return &pb.SendMsgResponse{IsOk: false}, fmt.Errorf("connList is nil")
	}
	log.Print(len(s.connList))
	for _, conn := range s.connList {
		if conn.userId == in.TextMsg.Receiver.UserId {
			log.Print("arrive if")
			err := conn.stream.SendMsg(in.GetTextMsg())
			log.Print("arrive after err")
			if err != nil {
				grpclog.Errorf("Error with stream %v. Error: %v", conn.stream, err)
				conn.err <- err
				return &pb.SendMsgResponse{IsOk: false}, err
			}
			break
		}
	}
	return &pb.SendMsgResponse{IsOk: true}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &server{})
	pb.RegisterCommunicationServer(s, &commServer{connList: make([]*Connection, 0)})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
