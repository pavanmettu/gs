package simplegossip

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	pb "github.com/pmettu/gs/pkg/proto/simplegossip"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

var (
	portid   = flag.Int("port", 10000, "Server Port")
	numnodes = flag.Int("numnodes", 16, "Total Number of Nodes")
	nodeid   = flag.Int("nodeid", -1, "Node ID")
)

type gossipServer struct {
	pb.UnimplementedGossipServiceServer
	pm       sync.Mutex
	nodeid   int32
	portid   int32
	numnodes int32
}

type gossipTuple struct {
	gmsg  string
	gpath []int32
}

var gc = map[string]gossipTuple{}

// CLIENT: Gossip to these nodes
func gossipnodes(s *gossipServer, gt gossipTuple, nodes []int32) {
	for i := 0; i < len(nodes); i++ {
		var gconn *grpc.ClientConn

		pt := s.portid + nodes[i]
		nodeaddr := fmt.Sprintf("localhost:%d", pt)
		fmt.Println(nodeaddr)
		gconn, err := grpc.Dial(nodeaddr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Cannot connect to gRPC server %s\n", nodeaddr)
		}
		defer gconn.Close()

		gms := &GossipMessageStruct{Sendernodeid: s.nodeid, Rcvrnodeid: nodes[i], Gmessage: gt.gmsg, Nodepaths: gt.gpath}
		c := pb.NewGossipServiceClient(gconn)
		resp, err := c.GossipMessage(context.Background(), &gms)
		if err != nil {
			log.Fatalf("Error sending Gossip message...\n")
		}
		fmt.Printf("Response is %v\n", resp)
	}
}

/*
 * SubmitMessage: Submits a message to the network after writing in its own cache.
 */
func (s *gossipServer) SubmitMessage(ctx context.Context, msg *pb.SubmitMessageStruct) (*pb.SubmitMessageRes, error) {
	fmt.Println("in SubmitMessage..")
	// Hash the message first

	hmd5 := md5.Sum([]byte(msg.gmessage))
	val, ok := gc[string(hmd5[:])]
	if !ok {
		// Gossip only if we don't have the message in our system
		var gt gossipTuple
		gt.gmsg = msg.gmessage
		gt.gpath = append(gt.gpath, s.nodeid)
		gc[string(hmd5[:])] = gt
		// Now gossip to nodes depending on if nodeid is even or odd.
		var nodes []int32
		count := 1
		if s.nodeid%2 == 0 {
			// Calculate the nodes which need to be sent
			i := s.nodeid + 1
			if i > 16 {
				i = 0
			}
			for i <= 16 {
				nodes = append(nodes, i)
				i++
				if count <= 3 && i > 16 {
					i = 1
				}
				count++
			}
		} else {
			i := s.nodeid - 1
			if i < 1 {
				i = 16
			}
			for i >= 0 {
				nodes = append(nodes, i)
				i--
				if count <= 3 && i < 1 {
					i = 16
				}
				count++
			}
		}
		// Send gossip to all nodes connected to this node
		go gossipnodes(s, gt, nodes)
	}
	sres := &SubmitMessageRes{Messageadded: true, Messageid: string(hmd5[:])}
	return

}

func (s *gossipServer) QueryMessage(ctx context.Context, msg *pb.QueryMessageStruct) (*pb.QueryMessageRes, error) {
	fmt.Println("In QueryMessage..")
	var qres QueryMessageRes
	gmsg, ok := gc[msg.messageid]
	if ok {
		qres.Messagefound = true
		qres.Gmessage = gmsg
	} else {
		qres.Messagefound = false
	}
	return qres, nil
}

func (s *gossipServer) ListMessage(ctx context.Context, msg *pb.ListMessageStruct) (*pb.ListMessageRes, error) {
	fmt.Println("In ListMessage.")
	// Go through the list of messages and join them into one
	var lres ListMessageRes
	lres.Moremessages = false
	for u, v := range gc {
		lres.Gmessages = append(lres.Gmessages, v.gmsg)
	}
	return lres, nil
}

func (s *gossipServer) GossipMessage(ctx context.Context, msg *pb.GossipMessageStruct) (*pb.GossipMessageRes, error) {
	fmt.Println("In GossipMessage")
	hmd5 := md5.Sum([]byte(msg.gmessage))
	val, ok := gc[string(hmd5[:])]
	if !ok {
		// Gossip only if we don't have the message in our system
		var gt gossipTuple
		gt.gmsg = msg.gmessage
		for i := 0; i < len(msg.gpath); i++ {
			gt.gpath = append(gt.gpath, msg.nodepaths[i])
		}
		gt.gpath = append(gt.gpath, s.nodeid)
		gc[string(hmd5[:])] = gt
		// Now gossip to nodes depending on if nodeid is even or odd.
		var nodes []int32
		count := 1
		if s.nodeid%2 == 0 {
			// Calculate the nodes which need to be sent
			i := s.nodeid + 1
			if i > 16 {
				i = 0
			}
			for i <= 16 {
				nodes = append(nodes, i)
				i++
				if count <= 3 && i > 16 {
					i = 1
				}
				count++
			}
		} else {
			i := s.nodeid - 1
			if i < 1 {
				i = 16
			}
			for i >= 0 {
				nodes = append(nodes, i)
				i--
				if count <= 3 && i < 1 {
					i = 16
				}
				count++
			}
		}
		// Send gossip to all nodes connected to this node
		go gossipnodes(s, gt, nodes)
	}
	sres := &GossipMessageRes{Rcvrnodeid: msg.Rcvrnodeid, Msgaccepted: true}
	return sres, nil
}

func (s *gossipServer) ResyncMessages(ctx context.Context, msg *pb.ResyncMessagesStruct) (*pb.ResyncMessagesRes, error) {
	fmt.Println("In ResyncMessages..")
	return nil, nil
}

func main() {
	// Listen on port
	pt := *portid + *nodeid
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", pt))
	if err != nil {
		log.Fatalf("Failed to listen on port %d %v", *portid, err)
	}

	// Start GRPC Server
	s := gossipServer{}
	s.nodeid = int32(*nodeid)
	s.portid = int32(*portid)
	s.numnodes = int32(*numnodes)
	grpcServer := grpc.NewServer()
	pb.RegisterGossipServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start service on port %d", *portid)
	}
}
