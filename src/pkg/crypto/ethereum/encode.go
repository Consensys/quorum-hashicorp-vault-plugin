package ethereum

import (
	"github.com/consensys/quorum/common"
	"github.com/consensys/quorum/rlp"
	"golang.org/x/crypto/sha3"
)

func Hash(object interface{}) (hash common.Hash, err error) {
	hashAlgo := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hashAlgo, object)
	if err != nil {
		return common.Hash{}, err
	}
	hashAlgo.Sum(hash[:0])
	return hash, nil
}
