package utils

import (
	"github.com/johankristianss/etherspace/pkg/security/crypto"

	"github.com/google/uuid"
)

func GenerateRandomID() string {
	uuid := uuid.New()
	crypto := crypto.CreateCrypto()
	return crypto.GenerateHash(uuid.String())
}
