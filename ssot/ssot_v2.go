package ssot

import (
	"crypto/ed25519"
	"encoding/binary"

	"github.com/frankbraun/codechain/util/base64"
	"github.com/frankbraun/codechain/util/hex"
)

// SignedHeadV2 is a signed Codechain head ready for publication as a SSOT with
// DNS TXT records (version 2).
type SignedHeadV2 struct {
	version      uint8    // the version of the signed head
	pubKey       [32]byte // Ed25519 public key of SSOT head signer
	pubKeyRotate [32]byte // Ed25519 pubkey to rotate to, all 0 if unused
	validFrom    int64    // this signed head is valid from the given Unix time
	validTo      int64    // this signed head is valid to the given Unix time
	counter      uint64   // signature counter
	head         [32]byte // the Codechain head to sign
	line         uint32   // the last signed line number
	signature    [64]byte // signature with pubkey over all previous fields
}

// marshal signed head without signature.
func (sh *SignedHeadV2) marshal() [125]byte {
	var m [125]byte
	var b [8]byte
	var l [4]byte
	m[0] = sh.version
	copy(m[1:33], sh.pubKey[:])
	copy(m[33:65], sh.pubKeyRotate[:])
	binary.BigEndian.PutUint64(b[:], uint64(sh.validFrom))
	copy(m[65:73], b[:])
	binary.BigEndian.PutUint64(b[:], uint64(sh.validTo))
	copy(m[73:81], b[:])
	binary.BigEndian.PutUint64(b[:], sh.counter)
	copy(m[81:89], b[:])
	copy(m[89:121], sh.head[:])
	binary.BigEndian.PutUint32(l[:], sh.line)
	copy(m[121:125], l[:])
	return m
}

// Marshal signed head with signature and encode it as base64.
func (sh *SignedHeadV2) Marshal() string {
	var m [189]byte
	b := sh.marshal()
	copy(m[:125], b[:])
	copy(m[125:189], sh.signature[:])
	return base64.Encode(m[:])
}

func unmarshalV2(signedHead string) (*SignedHeadV2, error) {
	m, err := base64.Decode(signedHead, 189)
	if err != nil {
		return nil, err
	}
	var sh SignedHeadV2
	sh.version = m[0]
	copy(sh.pubKey[:], m[1:33])
	copy(sh.pubKeyRotate[:], m[33:65])
	sh.validFrom = int64(binary.BigEndian.Uint64(m[65:73]))
	sh.validTo = int64(binary.BigEndian.Uint64(m[73:81]))
	sh.counter = binary.BigEndian.Uint64(m[81:89])
	copy(sh.head[:], m[89:121])
	sh.line = binary.BigEndian.Uint32(m[121:125])
	copy(sh.signature[:], m[125:189])
	msg := sh.marshal()
	if !ed25519.Verify(sh.pubKey[:], msg[:], sh.signature[:]) {
		return nil, ErrSignedHeadSignature
	}
	return &sh, nil
}

// Version returns the version.
func (sh *SignedHeadV2) Version() int {
	return int(sh.version)
}

// Head returns the signed head.
func (sh *SignedHeadV2) Head() string {
	return hex.Encode(sh.head[:])
}

// PubKey returns the public key in base64 notation.
func (sh *SignedHeadV2) PubKey() string {
	return base64.Encode(sh.pubKey[:])
}

// PubKeyRotate returns the public key rotate in base64 notation.
func (sh *SignedHeadV2) PubKeyRotate() string {
	return base64.Encode(sh.pubKeyRotate[:])
}

// ValidFrom returns the valid from field of signed head.
func (sh *SignedHeadV2) ValidFrom() int64 {
	return sh.validFrom
}

// ValidTo returns the valid to field of signed head.
func (sh *SignedHeadV2) ValidTo() int64 {
	return sh.validTo
}

// Counter returns the counter of signed head.
func (sh *SignedHeadV2) Counter() uint64 {
	return sh.counter
}

// Line returns the last signed line number of signed head.
func (sh *SignedHeadV2) Line() int {
	return int(sh.line)
}

// Signature returns the base64-encoded signature of the signed head.
func (sh *SignedHeadV2) Signature() string {
	return base64.Encode(sh.signature[:])
}

// HeadBuf returns the signed head.
func (sh *SignedHeadV2) HeadBuf() [32]byte {
	var b [32]byte
	copy(b[:], sh.head[:])
	return b
}
