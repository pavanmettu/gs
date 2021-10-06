package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"pkg/proto/simplegossip"
	"time"

	"google.golang.org/grpc"
)

var (
	nodeid  = flag.Int("nodeid", 0, "NodeID to send/query to/from")
	portid  = flag.Int("portid", 0, "PortID of Server")
	message = flag.String("msg", "", "Message to send")
	qmsg    = flag.String("qmsg", "", "Message ID to query")
	cmd     = flag.String("command", "", "Command to send to server.")
)

func gssendmsg(msg *string, nodeid *int32, portid *int32) {
	var opts []grpc.DialOption
	pt := *portid + *nodeid
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := &simplegossip.NewGossipServiceClientClient(conn)
	newMsg := &simplegossip.SubmitMessageStruct{Nodeid: *nodeid, Gmessage: *msg}
	fmt.Printf("In gssendmsg. %s, %d, %d\n", *msg, *nodeid, *portid)
	res, err := client.SubmitMessage(ctx, newMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func gsquerymsg(qmsg *string, nodeid *int32, portid *int32) {
	var opts []grpc.DialOption
	pt := *portid + *nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := &simplegossip.NewGossipServiceClientClientClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	qMsg := &simplegossip.QueryMessageStruct{Nodeid: *nodeid, Messageid: *qmsg}
	res, err := client.SubmitMessage(ctx, qMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func glistmsg(nodeid *int32, portid *int32) {
	var opts []grpc.DialOption
	pt := *portid + *nodeid
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
	client := &simplegossip.NewGossipServiceClientClient(conn)
	lMsg := &simplegossip.ListMessageStruct{Nodeid: *nodeid, Nummsgs: 100}
	res, err := client.SubmitMessage(ctx, lMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func main() {

	flag.Parse()
	switch *cmd {
	case "SendMsg":
		gssendmsg(message, *nodeid, *portid)
	case "QueryMsg":
		gsquerymsg(qmsg, *nodeid, *portid)
	case "ListMsg":
		gslistmsg(nodeid, portid)
	case "default":
		log.Fatalf("No Command")

	}
}
