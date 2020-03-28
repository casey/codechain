package ssot

import (
	"errors"
)

// ErrSignedHeadSignature is returned if a signed head signature does not verify.
var ErrSignedHeadSignature = errors.New("ssot: signed head signature does not verify")

// ErrSignedHeadFuture is returned if the validity of a signed head is in the future.
var ErrSignedHeadFuture = errors.New("ssot: signed head is valid in the future")

// ErrSignedHeadExpired is returned if the validity of a signed head is expired.
var ErrSignedHeadExpired = errors.New("ssot: signed head is expired")

// ErrValidityTooLong is returned if the validity is too long.
var ErrValidityTooLong = errors.New("ssot: validity is too long")

// ErrValidityTooShort is returned if the validity is too short.
var ErrValidityTooShort = errors.New("ssot: validity is too short")

// ErrTXTNoValidHead is returned if no valid TXT record for head could be found.
var ErrTXTNoValidHead = errors.New("ssot: no valid TXT record for head found")

// ErrTXTNoValidURL is returned if no valid TXT record for URL could be found.
var ErrTXTNoValidURL = errors.New("ssot: no valid TXT record for URL found")
