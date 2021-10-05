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

func main() {

	var opts []grpc.DialOption
	flag.Parse()
	opts = append(opts, grpc.WithBlock())
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to Server Service..%v\n", err)
	}
	defer conn.Close()
	gclient := pb.NewGossipServiceClient()

	switch *cmd {
	case "SendMsg":
		gssendmsg(gclient, message, nodeid, portid)
	case "QueryMsg":
		gsquerymsg(client, qmsg, nodeid, portid)
	case "ListMsg":
		gslistmsg(client, nodeid, portid)
	case "default":
		log.Fatalf("No Command")

	}
}
