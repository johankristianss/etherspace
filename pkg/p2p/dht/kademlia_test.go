package dht

import (
	"context"
	"fmt"
	"testing"
	"time"

	net "github.com/johankristianss/evrium/pkg/p2p/network"
	"github.com/johankristianss/evrium/pkg/p2p/network/mock"
	"github.com/johankristianss/evrium/pkg/security/crypto"
	"github.com/stretchr/testify/assert"
)

func createKademliaNode(t *testing.T, n mock.Network, addr string) *Kademlia {
	crypto := crypto.CreateCrypto()
	prvKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)

	node := net.Node{Addr: addr}
	contact1, err := CreateContact(node, prvKey)
	assert.Nil(t, err)

	m := mock.CreateMessenger(n, node)
	k, err := CreateKademlia(m, contact1)
	assert.Nil(t, err)

	return k
}

func TestKademliaFindRemoteContacts(t *testing.T) {
	n := mock.CreateFakeNetwork()

	b := createKademliaNode(t, n, "localhost:8000")
	k := createKademliaNode(t, n, "localhost:8001")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	k.ping(b.Contact.Node, ctx)
	cancel()

	targetID := k.Contact.ID

	// Ask bootstrap node for contact info to the k node
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	contacts, err := k.findRemoteContacts(b.Contact.Node, targetID.String(), 100, ctx)
	cancel()
	assert.Nil(t, err)

	foundContact := false
	for _, contact := range contacts {
		if contact.ID.String() == targetID.String() {
			foundContact = true
		}
	}

	assert.True(t, foundContact)

	b.Shutdown()
	k.Shutdown()
}

func TestKademliaFindContacts(t *testing.T) {
	n := mock.CreateFakeNetwork()

	bootstrapNode := createKademliaNode(t, n, "localhost:8000")

	var nodes []*Kademlia
	for i := 1; i < 200; i++ {
		k := createKademliaNode(t, n, "localhost:800"+fmt.Sprint(i))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		k.RegisterNetwork(bootstrapNode.Contact.Node, k.Contact.ID.String(), ctx)
		cancel()
		nodes = append(nodes, k)
	}

	targetID := nodes[10].Contact.ID
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	contacts, err := nodes[4].FindContacts(targetID.String(), 10, ctx)
	cancel()
	assert.Nil(t, err)
	assert.True(t, len(contacts) > 0)
	assert.Equal(t, contacts[0].ID.String(), targetID.String())

	targetID = nodes[3].Contact.ID
	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	contacts, err = nodes[10].FindContacts(targetID.String(), 10, ctx)
	cancel()
	assert.Nil(t, err)
	assert.True(t, len(contacts) > 0)
	assert.Equal(t, contacts[0].ID.String(), targetID.String())
}

func TestKademliaFindContact(t *testing.T) {
	n := mock.CreateFakeNetwork()

	bootstrapNode := createKademliaNode(t, n, "localhost:8000")

	var nodes []*Kademlia
	for i := 1; i < 20; i++ {
		k := createKademliaNode(t, n, "localhost:800"+fmt.Sprint(i))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		k.RegisterNetwork(bootstrapNode.Contact.Node, k.Contact.ID.String(), ctx)
		cancel()
		nodes = append(nodes, k)
	}

	targetID := nodes[10].Contact.ID
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	contact, err := nodes[0].FindContact(targetID.String(), ctx)
	cancel()
	assert.Nil(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, contact.ID.String(), targetID.String())
}

func TestKademliaPutGetRemote(t *testing.T) {
	n := mock.CreateFakeNetwork()

	bootstrapNode := createKademliaNode(t, n, "localhost:8000")

	var nodes []*Kademlia
	for i := 1; i < 20; i++ {
		k := createKademliaNode(t, n, "localhost:800"+fmt.Sprint(i))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		k.RegisterNetwork(bootstrapNode.Contact.Node, k.Contact.ID.String(), ctx)
		cancel()
		nodes = append(nodes, k)
	}

	crypto := crypto.CreateCrypto()
	prvKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)
	id, err := crypto.GenerateID(prvKey)
	assert.Nil(t, err)

	targetNode := nodes[10]
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = nodes[5].putRemote(targetNode.Contact.Node, id, prvKey, "/key1", "test1", ctx)
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	kvs, err := nodes[6].getRemote(targetNode.Contact.Node, id, "/key1", ctx)
	cancel()
	assert.Nil(t, err)
	assert.Equal(t, len(kvs), 1)
	assert.Equal(t, kvs[0].Value, "test1")

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = nodes[5].putRemote(targetNode.Contact.Node, id, prvKey, "/key2", "test2", ctx)
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	kvs, err = nodes[6].getRemote(targetNode.Contact.Node, id, "", ctx)
	cancel()
	assert.Nil(t, err)
	assert.Equal(t, len(kvs), 2)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	_, err = nodes[6].getRemote(targetNode.Contact.Node, "/prefix/not_found", "", ctx)
	cancel()
	assert.NotNil(t, err)
}

func TestKademliaPutGet(t *testing.T) {
	n := mock.CreateFakeNetwork()

	bootstrapNode := createKademliaNode(t, n, "localhost:8000")

	var nodes []*Kademlia
	for i := 1; i < 50; i++ {
		k := createKademliaNode(t, n, "localhost:800"+fmt.Sprint(i))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		k.RegisterNetwork(bootstrapNode.Contact.Node, k.Contact.ID.String(), ctx)
		cancel()
		nodes = append(nodes, k)
	}

	replicationFactor := 5

	crypto := crypto.CreateCrypto()
	prvKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)
	id, err := crypto.GenerateID(prvKey)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = nodes[5].Put(id, prvKey, "/prefix/key1", "test1", replicationFactor, ctx)
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	kvs, err := nodes[28].Get(id, "/prefix/key1", replicationFactor, ctx)
	cancel()
	assert.Nil(t, err)
	assert.Equal(t, len(kvs), 1)
	assert.Equal(t, kvs[0].Value, "test1")

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = nodes[12].Put(id, prvKey, "/prefix/key2", "test2", replicationFactor, ctx)
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	kvs, err = nodes[40].Get(id, "/prefix", replicationFactor, ctx)
	cancel()
	assert.Nil(t, err)

	count := 0
	for _, kv := range kvs {
		if kv.Value == "test1" || kv.Value == "test2" {
			count++
		}
	}
	assert.Equal(t, count, 2)
}

func TestKademliaNodeRegAndLookup(t *testing.T) {
	n := mock.CreateFakeNetwork()

	bootstrapNode := createKademliaNode(t, n, "localhost:8000")

	var nodes []*Kademlia
	for i := 1; i < 20; i++ {
		k := createKademliaNode(t, n, "localhost:800"+fmt.Sprint(i))
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		k.RegisterNetwork(bootstrapNode.Contact.Node, k.Contact.ID.String(), ctx)
		cancel()
		nodes = append(nodes, k)
	}

	crypto := crypto.CreateCrypto()
	prvKey, err := crypto.GeneratePrivateKey()
	assert.Nil(t, err)
	id, err := crypto.GenerateID(prvKey)
	assert.Nil(t, err)

	node := net.CreateNode("testnodename", "192.168.1.1")
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = nodes[5].RegisterNode(id, prvKey, &node, ctx)
	cancel()
	assert.Nil(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
	lookupedNode, err := nodes[10].LookupNode(id, "testnodename", ctx)
	cancel()
	assert.Nil(t, err)
	assert.True(t, lookupedNode.Equals(node))
}
