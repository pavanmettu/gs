package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "google.golang.org/gs/simplegossip"

	"google.golang.org/grpc"
)

var (
	nodeid  = flag.Int("nodeid", 0, "NodeID to send/query to/from")
	portid  = flag.Int("portid", 0, "PortID of Server")
	message = flag.String("msg", "", "Message to send")
	qmsg    = flag.String("qmsg", "", "Message ID to query")
	cmd     = flag.String("cmd", "", "Command to send to server.")
)

func gssendmsg(msg string, nodeid int, portid int) {
	var opts []grpc.DialOption
	pt := portid + nodeid
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	fmt.Printf("In gssendmsg 1 %d %s\n", pt, serverAddr)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := pb.NewGossipServiceClient(conn)
	newMsg := &pb.SubmitMessageStruct{Nodeid: int(nodeid), Gmessage: msg}
	fmt.Printf("In gssendmsg. %v\n", newMsg)
	res, err := client.SubmitMessage(ctx, newMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func gsquerymsg(qmsg string, nodeid int, portid int) {
	var opts []grpc.DialOption
	pt := portid + nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := pb.NewGossipServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	qMsg := &pb.QueryMessageStruct{Nodeid: int(nodeid), Messageid: qmsg}
	res, err := client.QueryMessage(ctx, qMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func gslistmsg(nodeid int, portid int) {
	var opts []grpc.DialOption
	pt := portid + nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client := pb.NewGossipServiceClient(conn)
	lMsg := &pb.ListMessageStruct{Nodeid: int(nodeid), Nummsgs: 100}
	res, err := client.ListMessage(ctx, lMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func main() {

	flag.Parse()
	switch *cmd {
	case "SubmitMsg":
		gssendmsg(*message, *nodeid, *portid)
	case "QueryMsg":
		gsquerymsg(*qmsg, *nodeid, *portid)
	case "ListMsg":
		gslistmsg(*nodeid, *portid)
	case "default":
		log.Fatalf("No Command")

	}
}
