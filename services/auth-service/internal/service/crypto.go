package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

// Argon2Hasher implements a real Argon2ID password hash and verification routine.
type Argon2Hasher struct {
	memory     uint32
	iterations uint32
	parallel   uint8
	keyLen     uint32
}

// NewArgon2Hasher configures a hashing instance using Argon2ID.
func NewArgon2Hasher(memory uint32, iterations uint32, parallelism uint8) *Argon2Hasher {
	if memory == 0 {
		memory = 64 * 1024
	}
	if iterations == 0 {
		iterations = 3
	}
	if parallelism == 0 {
		parallelism = 2
	}
	return &Argon2Hasher{memory: memory, iterations: iterations, parallel: parallelism, keyLen: 32}
}

func (h *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, h.iterations, h.memory, h.parallel, h.keyLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		h.memory,
		h.iterations,
		h.parallel,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func safeUint32(value int) (uint32, bool) {
	if value < 0 || value > int(^uint32(0)) {
		return 0, false
	}
	return uint32(value), true //nolint:gosec // validated against the uint32 range before conversion.
}

func (h *Argon2Hasher) Compare(storedHash, password string) bool {
	parts := strings.Split(strings.TrimPrefix(storedHash, "$"), "$")
	if len(parts) < 5 {
		return false
	}
	if parts[0] != "argon2id" {
		return false
	}
	params := strings.Split(parts[2], ",")
	var memory, iterations uint64
	var parallel uint64
	for _, part := range params {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return false
		}
		value, err := strconv.ParseUint(kv[1], 10, 64)
		if err != nil {
			return false
		}
		switch kv[0] {
		case "m":
			memory = value
		case "t":
			iterations = value
		case "p":
			parallel = value
		}
	}
	if memory == 0 || iterations == 0 || parallel == 0 {
		return false
	}
	if memory > uint64(^uint32(0)) || iterations > uint64(^uint32(0)) || parallel > uint64(^uint8(0)) {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	wantLen, ok := safeUint32(len(want))
	if !ok {
		return false
	}

	got := argon2.IDKey([]byte(password), salt, uint32(iterations), uint32(memory), uint8(parallel), wantLen)
	return subtle.ConstantTimeCompare(got, want) == 1
}

// StaticTokenIssuer produces deterministic opaque tokens appropriate for the scaffolded auth-service.
type StaticTokenIssuer struct{}

// NewStaticTokenIssuer constructs a token issuer suitable for local development.
func NewStaticTokenIssuer() *StaticTokenIssuer {
	return &StaticTokenIssuer{}
}

func (s *StaticTokenIssuer) IssueAccess(userID string, now time.Time) (string, error) {
	return fmt.Sprintf("access.%s.%d", userID, now.UnixNano()), nil
}

func (s *StaticTokenIssuer) IssueRefresh(userID, familyID string, now time.Time) (string, error) {
	return fmt.Sprintf("refresh.%s.%s.%d", userID, familyID, now.UnixNano()), nil
}
