package bip32

import (
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

const (
	// HardenedKeyStart is the index at which a hardened key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardened child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31
)

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.  See the package overview documentation for
// more details on how to use extended keys.
type ExtendedKey struct {
	key       []byte // This will be the pubkey for extended pub keys
	pubKey    []byte // This will only be set for extended priv keys
	chainCode []byte
	depth     uint8
	parentFP  []byte
	childNum  uint32
	version   []byte
	isPrivate bool
}

func Xkey(key string, args ...uint32) (*ExtendedKey, error) {
	c, e := newKeyFromString(key)
	if e != nil {
		return nil, e
	}
	for _, v := range args {
		c, e = c.derive(v)
		if e != nil {
			return nil, e
		}
	}
	return c, nil
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the provided seed
// will derive to an unusable secret key.  The ErrUnusable error will be
// returned if this should occur, so the caller must check for it and generate a
// new seed accordingly.
func NewMaster(seed []byte, HDPrivateKeyID [4]byte) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [minSeedBytes, maxSeedBytes].
	if len(seed) < minSeedBytes || len(seed) > maxSeedBytes {
		return nil, fmt.Errorf("seed length must be between %d and %d bits", minSeedBytes*8, maxSeedBytes*8)
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	secretKey := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Ensure the key in usable.
	secretKeyNum := new(big.Int).SetBytes(secretKey)
	if secretKeyNum.Cmp(btcec.S256().N) >= 0 || secretKeyNum.Sign() == 0 {
		return nil, errors.New("unusable seed")
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}
	return newExtendedKey(HDPrivateKeyID[:], secretKey, chainCode, parentFP, 0, 0, true), nil
}
