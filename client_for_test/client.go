package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	pb "gochat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultUserId = "jack"
)

var (
	addr     = flag.String("addr", "localhost:8888", "the address to connect to")
	userId   = flag.String("sender", defaultUserId, "userId must be set")
	receiver = flag.String("receiver", defaultUserId, "who receives msgs")
)

var wg sync.WaitGroup

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetIMToken(ctx, &pb.AuthRequest{UserId: *userId})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetImToken())
	commClient := pb.NewCommunicationClient(conn)
	myself := pb.User{
		UserId:   *userId,
		UserName: *userId + "张世豪",
		ImToken:  r.GetImToken(),
	}
	stream, err := commClient.ConnectServer(context.Background(), &myself)
	if err != nil {
		fmt.Errorf("Connect failed: %v", err)
		return
	}
	wg = sync.WaitGroup{}
	wg.Add(1)
	go ReceiveMsg(stream)
	wg.Add(1)
	go SendMsg(commClient, &myself, &pb.User{UserId: *receiver})
	wg.Wait()
}

func ReceiveMsg(commClient pb.Communication_ConnectServerClient) {
	defer wg.Done()
	for {
		msg, err := commClient.Recv()

		if err != nil {
			fmt.Errorf("Error reading message: %v", err)
			break
		}

		fmt.Printf("%v: %s\n", msg.Sender.GetUserName(), msg.Text)
	}
}

func SendMsg(commClient pb.CommunicationClient, sender *pb.User, receiver *pb.User) {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	ts := time.Now()
	msgID := sha256.Sum256([]byte(ts.String() + sender.UserId))
	for scanner.Scan() {
		msg := &pb.TextMsg{
			MsgId:     hex.EncodeToString(msgID[:]),
			Sender:    sender,
			Receiver:  receiver,
			Text:      scanner.Text(),
			Timestamp: ts.String(),
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		_, err := commClient.SendMsg(ctx, &pb.SendMsgRequest{TextMsg: msg, IsGroupMsg: false})
		if err != nil {
			fmt.Printf("Error sending message: %v", err)
			break
		}
	}
}
