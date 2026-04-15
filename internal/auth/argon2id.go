package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

// Argon2id parameters (OWASP-aligned; memory in KiB per golang.org/x/crypto/argon2).
const (
	argon2Time    = 3
	argon2Memory  = 65536 // 64 MiB
	argon2Threads = 4
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

// HashPassword returns an Argon2id hash of the password (PHC-style string).
func HashPassword(password string) ([]byte, error) {
	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	key := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, uint8(argon2Threads), argon2KeyLen)
	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64key := base64.RawStdEncoding.EncodeToString(key)
	// PHC-style: $argon2id$v=19$m=65536,t=3,p=4$<salt>$<hash>
	s := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		argon2Memory, argon2Time, argon2Threads, b64salt, b64key)
	return []byte(s), nil
}

func checkArgon2id(encoded, password string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return errors.New("invalid argon2id hash format")
	}
	if parts[2] != "v=19" {
		return errors.New("unsupported argon2id version")
	}
	mem, time, threads, err := parseArgon2Params(parts[3])
	if err != nil {
		return err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return errors.New("invalid argon2id salt")
	}
	want, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return errors.New("invalid argon2id digest")
	}
	got := argon2.IDKey([]byte(password), salt, time, mem, threads, uint32(len(want)))
	if subtle.ConstantTimeCompare(got, want) != 1 {
		return bcrypt.ErrMismatchedHashAndPassword
	}
	return nil
}

func parseArgon2Params(s string) (memory, time uint32, threads uint8, err error) {
	// m=65536,t=3,p=4
	for _, kv := range strings.Split(s, ",") {
		kv = strings.TrimSpace(kv)
		switch {
		case strings.HasPrefix(kv, "m="):
			v, e := strconv.ParseUint(strings.TrimPrefix(kv, "m="), 10, 32)
			if e != nil {
				return 0, 0, 0, fmt.Errorf("argon2id params: %w", e)
			}
			memory = uint32(v)
		case strings.HasPrefix(kv, "t="):
			v, e := strconv.ParseUint(strings.TrimPrefix(kv, "t="), 10, 32)
			if e != nil {
				return 0, 0, 0, fmt.Errorf("argon2id params: %w", e)
			}
			time = uint32(v)
		case strings.HasPrefix(kv, "p="):
			v, e := strconv.ParseUint(strings.TrimPrefix(kv, "p="), 10, 8)
			if e != nil {
				return 0, 0, 0, fmt.Errorf("argon2id params: %w", e)
			}
			threads = uint8(v)
		}
	}
	if memory == 0 || time == 0 || threads == 0 {
		return 0, 0, 0, errors.New("argon2id params: missing m, t, or p")
	}
	return memory, time, threads, nil
}
