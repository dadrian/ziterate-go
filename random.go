package ziterate

import (
	"io"
	mathrand "math/rand"
)

// This file only exists because we don't want to import "crypto/rand" and
// "math/rand" in another file possibly, because math/rand is usually the devil.

// NewSeedReader returns a deterministic random reader for repeatable iteration.
func NewSeedReader(seed uint64) io.Reader {
	return mathrand.New(mathrand.NewSource(int64(seed)))
}
