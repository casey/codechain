package ssot

import (
	"time"
)

// Valid checks if the signed head sh is currently valid
// (as defined by validFrom and validTo).
// It returns nil, if the signed check is valid and an error otherwise.
func Valid(sh SignedHead) error {
	now := time.Now().UTC().Unix()
	if now < sh.ValidFrom() {
		return ErrSignedHeadFuture
	}
	if now > sh.ValidTo() {
		return ErrSignedHeadExpired
	}
	return nil
}
