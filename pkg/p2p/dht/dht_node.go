package dht

import (
	"github.com/johankristianss/evrium/internal/crypto"
	"github.com/johankristianss/evrium/pkg/p2p/network/libp2p"
)

func CreateDHT(port int, name string) (DHT, error) {
	m, err := libp2p.CreateMessenger(port, name)
	if err != nil {
		return nil, err
	}

	id, err := crypto.CreateIdendity()
	if err != nil {
		return nil, err
	}

	prvKey := id.PrivateKeyAsHex()

	contact, err := CreateContact(m.Node, prvKey)
	if err != nil {
		return nil, err
	}

	k, err := CreateKademlia(m, contact)
	if err != nil {
		return nil, err
	}

	return k, nil
}
