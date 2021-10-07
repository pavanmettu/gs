```
# Simple Gossip Protocol
```
## Requirements:

1. Upto 16 nodes support.
2. Each node can communicate with 3 other nodes.
3. Supports two API's
    a. submit_message(message)
    b. get_messages() --> List[str]
4. All nodes eventually should get messages.
5. Network Calls(TCP/UDP/HTTP).
6. Message ordering not guarenteed.

## Non-Requirements but considered:

1. In-Memory Cache or Persistent File Cache.
2. Re-connection if a node goes down and comes back.
3. Static/Dynamic node reachability map(3 Nodes reachability).
4. Re-sync of messages already in cache if a node restarts.
5. Communication Security between Nodes.(no Authentication/TLS - For now Simple
Request/Response)
6. Dead Nodes.
7. Dockerising servers.

## Design:

1. Each node has an ID/Port and an Static/Dynamic Graph on what nodes it can reach.
2. Each node is listening on a port and maintain a connection to each of the nodes it can
reach.
3. Each node maintains a cache of messages it receives with each message in this
format:
msgID: Hash of Message
Message: Message string
Time Recvd.: For re-sync operations.
Recvd. From: Array of Nodes
All the messages are maintained in a HashMap with msgID as the key and the
Message and Array of Nodes
tuple as Value.
4. Extension:
1. Make data persistent and use a file for storing messages instead of in-memory
cache.
2. Support re-sync operation.

```
## Node Communication:
We are going to use ProtoBuf for communication between nodes.
```

Each node comes up and Listens on a port and connects to three other nodes using a
go routine for each node.
Each node is either given a Graph or nodes to connect or :wq
If any of the nodes go down/doesn't come up, a go routine continiously tries to connect
to that node.
All nodes are maintained in a array of structs with NodeID as a reference and a state.

```
There are going to be four Protobuf structures for maintaining Gossip.
```
1. Submitting a Message.
2. Gossiping a message.
3. List of messages.
4. Query a message.
Extended:
1. Re-sync Messages.

```
## Submitting a Message:
A client sends a message with an ID and a message.
$ gclient --id 1 --message "Apple"
xHviKlH
```
This operation sends a message to the Node 1 and gets an ack with a string. Eventually
all
the nodes gets this message through our gossip protocol. Node1 in this case goes
through its list
of nodes it has to send and sends to all nodes adding its ID to the message.

For example if Node1 has Node2, Node3, Node 4 as its neighbors, it sends a messages
to each of them
in this format.

```
Message, 1 (2)
Message, 1 (3)
Message, 1 (4)
```
Node 2 upon receiving the message will check if it has already seen it. If it finds, it takes
no action.
It not, it stores and sends the message to all its know Neighors in this format.

```
Message, 1,2 (5)
Message 1,2 (6)
Message 1,2 (7)
```

This way everyone sends to all their neighbors if they have not seen the message
adding their own ID to
the message.

## Quering a message:
A client upon receiving a response from submission can check if that message exists on
a node.
$ gclient --id 5 --msgID xHviKlH
Apple
Query is checking if a message exists on a server based on MessageID.
Check the HashMap for the msgID and if it exists return the message to client.

## Listing Messages
A client can send a list API to any of the servers and receive in return all messages the
server has.
$ gclient --id 4 --list
Apple
....
....

```
Go through the list of all messages on server and return the list to the client.
```
Re-Sync:
This is an in-memory-based Gossip Message System but nodes can go up and down.
We need to make sure all messages are synced when a node restarts. Considering
that there is no guarentee for messages to be stored in sequence, we have to
read all messages and pack them in a string buffer and send back to client.
To do that, we need to maintain a data structure of nodes from which messages are
received and connect to any of the nodes to sync messages.

Analytics:
TBD

Testing:
TBD
