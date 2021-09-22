package crypto

import (
	"bytes"
	"crypto/rand"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func NewBabyjubjub() (babyjubjub.PrivateKey, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return babyjubjub.PrivateKey{}, err
	}

	// Usually standards implementations of eddsa do not require the choice of a specific hash function (usually it's SHA256).
	// Here we needed to allow the choice of the hash, so we can choose a hash function that is easily programmable in a snark circuit.
	// Same hFunc should be used for sign and verify
	return babyjubjub.GenerateKey(bytes.NewReader(seed))
}
