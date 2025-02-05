package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"

	pb "google.golang.org/gs/simplegossip"

	"google.golang.org/grpc"
)

var (
	portid   = flag.Int("port", 10000, "Server Port")
	numnodes = flag.Int("numnodes", 4, "Total Number of Nodes")
	nodeid   = flag.Int("nodeid", -1, "Node ID")
)

type gossipServer struct {
	pb.UnimplementedGossipServiceServer
	pm       *sync.Mutex
	nodeid   int
	portid   int
	numnodes int
}

type gossipTuple struct {
	gmsg  string
	gpath []int
}

var gc = map[string]gossipTuple{}
var rf = map[int]int{}

// CLIENT: Gossip to these nodes
func gossipnodes(s *gossipServer, gt gossipTuple, nodes []int) {
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
		var gpth []int32
		for i := 0; i < len(gt.gpath); i++ {
			gpth = append(gpth, int32(gt.gpath[i]))
		}

		gms := &pb.GossipMessageStruct{Sendernodeid: int32(s.nodeid), Rcvrnodeid: int32(nodes[i]), Gmessage: gt.gmsg, Nodepaths: gpth}
		c := pb.NewGossipServiceClient(gconn)
		_, err = c.GossipMessage(context.Background(), gms)
		if err != nil {
			log.Fatalf("Error sending Gossip message...\n")
		}
	}
}

func getnodes(s *gossipServer, count int) []int {
	var nodes []int
	if s.nodeid%2 == 0 {
		// Calculate the nodes which need to be sent
		ncnt := 0
		i := s.nodeid + 1
		if i > s.numnodes {
			i = 1
		}
		for i <= s.numnodes {
			nodes = append(nodes, i)
			i++
			ncnt++
			if ncnt < count && i > s.numnodes {
				i = 1
			}
			if ncnt == count {
				break
			}
		}
	} else {
		ncnt := 0
		i := s.nodeid - 1
		if i < 1 {
			i = s.numnodes
		}
		for i >= 1 {
			nodes = append(nodes, i)
			i--
			ncnt++
			if ncnt == count {
				break
			}
			if ncnt < count && i < 1 {
				i = s.numnodes
			}
		}
	}
	return nodes
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

/*
 * SubmitMessage: Submits a message to the network after writing in its own cache.
 */
func (s *gossipServer) SubmitMessage(ctx context.Context, msg *pb.SubmitMessageStruct) (*pb.SubmitMessageRes, error) {
	// Hash the message first

	hmd5 := GetMD5Hash(msg.Gmessage)
	s.pm.Lock()
	_, ok := gc[hmd5]
	s.pm.Unlock()
	if !ok {
		// Gossip only if we don't have the message in our system
		var gt gossipTuple
		gt.gmsg = msg.Gmessage
		gt.gpath = append(gt.gpath, s.nodeid)
		s.pm.Lock()
		gc[hmd5] = gt
		s.pm.Unlock()
		// Now gossip to nodes depending on if nodeid is even or odd.
		var nodes []int
		count := 2
		nodes = getnodes(s, count)
		// Send gossip to all nodes connected to this node
		gossipnodes(s, gt, nodes)
	}
	sres := &pb.SubmitMessageRes{Messageadded: true, Messageid: hmd5}
	return sres, nil
}

func (s *gossipServer) QueryMessage(ctx context.Context, msg *pb.QueryMessageStruct) (*pb.QueryMessageRes, error) {
	fmt.Println("In QueryMessage..")
	qres := new(pb.QueryMessageRes)
	qmsg := new(pb.MsgFromNode)
	s.pm.Lock()
	gmsg, ok := gc[msg.Messageid]
	s.pm.Unlock()
	if ok {
		qres.Messagefound = true
		var gpth []int32
		for i := 0; i < len(gmsg.gpath); i++ {
			gpth = append(gpth, int32(gmsg.gpath[i]))
		}
		qmsg.Msg = gmsg.gmsg
		qmsg.Nodepath = gpth
		qres.Gmessage = qmsg
	} else {
		qres.Messagefound = false
	}
	return qres, nil
}

func (s *gossipServer) ListMessage(ctx context.Context, msg *pb.ListMessageStruct) (*pb.ListMessageRes, error) {
	fmt.Println("In ListMessage.")
	// Go through the list of messages and join them into one
	lres := new(pb.ListMessageRes)
	lres.Moremessages = false
	var res []string
	s.pm.Lock()
	for _, vmsg := range gc {
		res = append(res, vmsg.gmsg)
	}
	lres.Gmessages = res
	s.pm.Unlock()
	return lres, nil
}

func (s *gossipServer) GossipMessage(ctx context.Context, msg *pb.GossipMessageStruct) (*pb.GossipMessageRes, error) {
	fmt.Println("In GossipMessage")
	hmd5 := GetMD5Hash(msg.Gmessage)
	s.pm.Lock()
	_, ok := gc[hmd5]
	rf[int(msg.Rcvrnodeid)] += 1
	s.pm.Unlock()
	if !ok {
		// Gossip only if we don't have the message in our system
		var gt gossipTuple
		gt.gmsg = msg.Gmessage
		for i := 0; i < len(msg.Nodepaths); i++ {
			gt.gpath = append(gt.gpath, int(msg.Nodepaths[i]))
		}
		gt.gpath = append(gt.gpath, s.nodeid)
		s.pm.Lock()
		gc[hmd5] = gt
		s.pm.Unlock()
		// Now gossip to nodes depending on if nodeid is even or odd.
		var nodes []int
		count := 2
		nodes = getnodes(s, count)
		// Send gossip to all nodes connected to this node
		go gossipnodes(s, gt, nodes)
	}
	sres := &pb.GossipMessageRes{Rcvrnodeid: msg.Rcvrnodeid, Msgaccepted: true}
	return sres, nil
}

func (s *gossipServer) ResyncMessages(ctx context.Context, msg *pb.ResyncMessagesStruct) (*pb.ResyncMessagesRes, error) {
	fmt.Println("In ResyncMessages..")
	return nil, nil
}

func main() {
	// Listen on port
	flag.Parse()
	pt := *portid + *nodeid
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", pt))
	if err != nil {
		log.Fatalf("Failed to listen on port %d %v", *portid, err)
	}

	// Start GRPC Server
	s := gossipServer{}
	s.nodeid = int(*nodeid)
	s.portid = int(*portid)
	s.numnodes = int(*numnodes)
	s.pm = &sync.Mutex{}
	grpcServer := grpc.NewServer()
	pb.RegisterGossipServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start service on port %d", *portid)
	}
}
