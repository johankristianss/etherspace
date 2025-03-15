package network

type Message struct {
	ID      string
	From    Node
	To      Node
	Type    int
	Payload []byte
}
