package crypto

import "sync"

// PasswordEncoder interface for encoding passwords.
type PasswordEncoder interface {
	Gensalt() (salt string, err error)

	/**
	 * Encode the raw password. Generally, a good encoding algorithm applies a SHA-1 or
	 * greater hash combined with an 8-byte or greater randomly generated salt.
	 */
	Encode(password string, salt ...string) (encodedPassword string, err error)
	/**
	 * Matches verify the encoded password obtained from storage matches the submitted raw
	 * password after it too is encoded. Returns true if the passwords match, false if
	 * they do not. The stored password itself is never decoded.
	 *
	 * @param password the raw password to encode and match
	 * @param salt the encoded password or encodedPassword from storage to compare with
	 * @return true if the raw password, after encoding, matches the encoded password from
	 * storage
	 */
	Matches(password, salt string) (ok bool, err error)
}

var (
	encoderMu sync.RWMutex
	encoders  = make(map[string]PasswordEncoder)
)

// RegisterPasswordEncoder makes a PasswordEncoder available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterPasswordEncoder(name string, encoder PasswordEncoder) {
	encoderMu.Lock()
	defer encoderMu.Unlock()
	if encoder == nil {
		panic("crypto: Register PasswordEncoder is nil")
	}
	if _, dup := encoders[name]; dup {
		panic("crypto: Register called twice for PasswordEncoder " + name)
	}
	encoders[name] = encoder
}

// GetPasswordEncoder returns a register PasswordEncoder
func GetPasswordEncoder(name string) (encoder PasswordEncoder, ok bool) {
	encoder, ok = encoders[name]
	return
}
