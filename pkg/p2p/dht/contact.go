package dht

import (
	"fmt"
	"sort"

	"github.com/johankristianss/evrium/internal/crypto"
	net "github.com/johankristianss/evrium/pkg/p2p/network"
)

type Contact struct {
	ID       KademliaID `json:"kademliaid"`
	Node     net.Node   `json:"node"`
	prvKey   string
	distance KademliaID
}

func createContactWithKademliaID(id KademliaID, addr string) Contact { // Just for testing
	return Contact{ID: id, Node: net.Node{Addr: addr}}
}

func CreateContact(node net.Node, prvKey string) (Contact, error) {
	id, err := crypto.CreateIdendityFromString(prvKey)
	if err != nil {
		return Contact{}, err
	}

	return Contact{ID: CreateKademliaID(id.ID()), Node: node, prvKey: prvKey}, nil
}

func (contact *Contact) CalcDistance(target KademliaID) {
	contact.distance = contact.ID.CalcDistance(target)
}

func (contact *Contact) Less(otherContact Contact) bool {
	return contact.distance.Less(otherContact.distance)
}

func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Node.String())
}

type ContactCandidates struct {
	contacts []Contact
}

func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(candidates.contacts[j])
}
