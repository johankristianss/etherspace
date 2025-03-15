package crypto

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveredID(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	msg := "test"
	hash := GenerateHashFromString(msg)

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	recoveredID, err := RecoveredID(hash, []byte(string(signatureBytes)+"too_large_signature"))
	assert.NotNil(t, err)

	recoveredID, err = RecoveredID(hash, signatureBytes)
	assert.Nil(t, err)
	assert.Equal(t, idendity.ID(), recoveredID)
}

func TestRecoverFromStrings(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	msg := "test"
	hash := GenerateHashFromString(msg)

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	hash2, err := CreateHashFromString(hash.String())
	assert.Nil(t, err)

	recoveredID, err := RecoveredID(hash2, signatureBytes)
	assert.Nil(t, err)
	assert.Equal(t, idendity.ID(), recoveredID)
}

func TestRecoverFromStringsInvalidHash(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	msg := "test"
	hash := GenerateHashFromString(msg)

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	hash = GenerateHash([]byte("blablabla"))

	recoveredID, err := RecoveredID(hash, signatureBytes)
	assert.Nil(t, err)
	assert.NotEqual(t, idendity.ID(), recoveredID)
}

func TestRecoverPublicKey(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	hash := GenerateHash([]byte("test"))

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	pub, err := RecoverPublicKey(hash, signatureBytes)
	assert.Nil(t, err)
	assert.Equal(t, idendity.PublicKeyAsHex(), hex.EncodeToString(pub))
}

func TestRecoverPublicKeyInvalidSignature(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	idendity2, err := CreateIdendity()
	assert.Nil(t, err)

	hash := GenerateHash([]byte("test"))

	signatureBytes, err := Sign(hash, idendity2.PrivateKey())
	assert.Nil(t, err)

	pub, err := RecoverPublicKey(hash, signatureBytes)
	assert.Nil(t, err)
	assert.NotEqual(t, idendity.PublicKeyAsHex(), hex.EncodeToString(pub))
}

func TestSignAndVerify(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	msg := "test"
	hash := GenerateHashFromString(msg)

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	decPub, err := hex.DecodeString(idendity.PublicKeyAsHex())
	assert.Nil(t, err)

	ok, err := Verify(decPub, hash, signatureBytes)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestSignAndVerifyInvalidPubKey(t *testing.T) {
	idendity, err := CreateIdendity()
	assert.Nil(t, err)

	idendity2, err := CreateIdendity()
	assert.Nil(t, err)

	msg := "test"
	hash := GenerateHashFromString(msg)

	signatureBytes, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)

	decPub, err := hex.DecodeString(idendity2.PublicKeyAsHex())
	assert.Nil(t, err)

	ok, err := Verify(decPub, hash, signatureBytes)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestInterop(t *testing.T) {
	prvKey := "d6eb959e9aec2e6fdc44b5862b269e987b8a4d6f2baca542d8acaa97ee5e74f6"
	idendity, err := CreateIdendityFromString(prvKey)
	assert.Nil(t, err)

	// Alice creates an idendity
	// Alice shares the public key with Bob
	fmt.Println("Alice created an idendity, and shared the ID=SHA-2(public) key with Bob")
	fmt.Println("prvkey: " + idendity.PrivateKeyAsHex())
	fmt.Println("pubkey: " + idendity.PublicKeyAsHex())
	fmt.Println("id: " + idendity.ID())

	fmt.Println("")

	msg := "hello"
	hash := GenerateHashFromString(msg)

	signature, err := Sign(hash, idendity.PrivateKey())
	assert.Nil(t, err)
	signatureStr := hex.EncodeToString(signature)

	fmt.Println("Alice sends a message to Bob")
	fmt.Println("message: " + msg)
	fmt.Println("digest: " + hash.String())
	fmt.Println("signature: " + string(signatureStr))

	fmt.Println("")
	fmt.Println("Bob receives the message and the signature")
	// Bob receives the signature and the message

	signatureHex := "e713a1bb015fecabb5a084b0fe6d6e7271fca6f79525a634183cfdb175fe69241f4da161779d8e6b761200e1cf93766010a19072fa778f9643363e2cfadd640900"
	signatureBytes, err := hex.DecodeString(signatureHex)
	recoveredID, err := RecoveredID(hash, signatureBytes)

	fmt.Println("Bob recovers the id from the signature")
	// Bob checks that the recovered id is the same as the one Alice shared
	fmt.Println("recovered id: " + recoveredID)
	assert.Nil(t, err)
	assert.Equal(t, recoveredID, idendity.ID())
	fmt.Println("If the recovered id is the same as the one Alice shared, the signature is valid, and the message is authentic")
}
