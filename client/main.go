package main

import (
	"context"
	pb "github.com/pmettu/gs/pkg/proto/simplegossip"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	nodeid  = flag.Int("nodeid", "0", "NodeID to send/query to/from")
	portid  = flag.Int("portid", "0", "PortID of Server")
	message = flag.String("msg", "", "Message to send")
	qmsg    = flag.String("qmsg", "", "Message ID to query")
	cmd     = flag.String("command", "", "Command to send to server.")
)

func gssendmsg(msg *string, nodeid *int32, portid *int32) {
	pt := *portid + *nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := pb.NewPolicyServiceClient(conn)
	newMsg := &pb.SubmitMessageStruct{Nodeid: *nodeid, Gmessage: *msg}
	fmt.Printf("In gssendmsg. %s, %d, %d\n", *msg, *nodeid, *portid)
	res, err := client.SubmitMessage(ctx, newMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func gsquerymsg(qmsg *string, nodeid *int32, portid *int32) {
	pt := *portid + *nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := pb.NewPolicyServiceClient(conn)

	qMsg := &pb.QueryMessageStruct{Nodeid: *nodeid, Messageid: *qmsg}
	res, err := client.SubmitMessage(ctx, qMsg)
	if err != nil {
		log.Fatalf("Error from client Submit Message. %v\n", err)
	}
	log.Println(res)
}

func glistmsg(nodeid *int32, portid *int32) {
	pt := *portid + *nodeid
	serverAddr := fmt.Sprintf("localhost:%d", pt)
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	client := pb.NewPolicyServiceClient(conn)
	lMsg := &pb.ListMessageStruct{Nodeid: *nodeid, Nummsgs: 100}
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
		gssendmsg(message, nodeid, portid)
	case "QueryMsg":
		gsquerymsg(qmsg, nodeid, portid)
	case "ListMsg":
		gslistmsg(nodeid, portid)
	case "default":
		log.Fatalf("No Command")

	}
}
