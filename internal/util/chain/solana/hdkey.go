package solana

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// HardenedKeyStart represents the start index for hardened keys
	HardenedKeyStart uint32 = 0x80000000
	// MasterSecret is the master secret for generating master keys
	MasterSecret = "ed25519 seed"
)

// HDKey represents a hierarchical deterministic key
type HDKey struct {
	Version           uint32
	Depth             uint8
	ParentFingerprint uint32
	ChildIndex        uint32
	ChainCode         []byte
	PrivateKey        []byte
}

// FromMasterSeed derives the master HDKey from a seed
func FromMasterSeed(seed []byte) (*HDKey, error) {
	// Check seed length (128-512 bits; 256 bits is advised)
	seedLen := len(seed)
	if seedLen*8 < 128 || seedLen*8 > 512 {
		return nil, fmt.Errorf("seed length must be between 128 and 512 bits; got %d bits", seedLen*8)
	}

	// Create HMAC-SHA512 from seed
	hmacObj := hmac.New(sha512.New, []byte(MasterSecret))
	hmacObj.Write(seed)
	I := hmacObj.Sum(nil)

	// Split into left 32 bytes (private key) and right 32 bytes (chain code)
	privKey := I[:32]
	chainCode := I[32:]

	// Create and return HDKey
	key := &HDKey{
		Depth:             0,
		ParentFingerprint: 0,
		ChildIndex:        0,
		ChainCode:         chainCode,
		PrivateKey:        privKey,
	}

	return key, nil
}

// PublicKey returns the public key for the HDKey
func (k *HDKey) PublicKey() []byte {
	pub := make([]byte, ed25519.PublicKeySize+1)
	pub[0] = 0x00 // Version byte
	copy(pub[1:], ed25519.PublicKey(k.PrivateKey))
	return pub
}

// Sign signs the message with the private key
func (k *HDKey) Sign(message []byte) ([]byte, error) {
	if len(message) != 32 {
		return nil, errors.New("message must be 32 bytes")
	}

	privateKey := ed25519.NewKeyFromSeed(k.PrivateKey)
	return ed25519.Sign(privateKey, message), nil
}

// Verify verifies the signature against the message using the public key
func (k *HDKey) Verify(message, signature []byte) bool {
	if len(message) != 32 || len(signature) != 64 {
		return false
	}

	publicKey := ed25519.PublicKey(k.PublicKey()[1:]) // Skip version byte
	return ed25519.Verify(publicKey, message, signature)
}

func (k *HDKey) Derive(path string) (*HDKey, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, fmt.Errorf("invalid path format")
	}

	parts := strings.Split(path[2:], "/")
	current := k

	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		// Remove hardened marker
		hardened := strings.HasSuffix(part, "'")
		if hardened {
			part = part[:len(part)-1]
		}

		index, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid index: %s", part)
		}

		if hardened {
			index += 0x80000000
		}

		current, err = current.DeriveChild(uint32(index))
		if err != nil {
			return nil, err
		}
	}

	return current, nil
}

// DeriveChild derives a child key from the current key
func (k *HDKey) DeriveChild(index uint32) (*HDKey, error) {
	// Only support hardened child keys for Ed25519
	if index < 0x80000000 {
		return nil, fmt.Errorf("non-hardened child derivation not possible for Ed25519")
	}

	data := make([]byte, 37)
	data[0] = 0x00
	copy(data[1:], k.PrivateKey)
	copy(data[33:], uint32ToBytes(index))

	hmacObj := hmac.New(sha512.New, k.ChainCode)
	hmacObj.Write(data)
	I := hmacObj.Sum(nil)

	child := &HDKey{
		Depth:             k.Depth + 1,
		ParentFingerprint: k.Fingerprint(),
		ChildIndex:        index,
		ChainCode:         I[32:],
		PrivateKey:        I[:32],
	}

	return child, nil
}

// Fingerprint returns the fingerprint of the key
func (k *HDKey) Fingerprint() uint32 {
	// For simplicity, using first 4 bytes of public key as fingerprint
	return uint32(k.PrivateKey[0])<<24 | uint32(k.PrivateKey[1])<<16 |
		uint32(k.PrivateKey[2])<<8 | uint32(k.PrivateKey[3])
}

func uint32ToBytes(i uint32) []byte {
	bytes := make([]byte, 4)
	bytes[0] = byte(i >> 24)
	bytes[1] = byte(i >> 16)
	bytes[2] = byte(i >> 8)
	bytes[3] = byte(i)
	return bytes
}
