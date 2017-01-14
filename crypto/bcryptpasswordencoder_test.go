package crypto

import "testing"

func TestBcryptPasswordEncoder(t *testing.T) {
	const passwd = "123456abc"

	encoder := &bcryptPasswordEncoder{}

	passwordHash, err := encoder.Encode(passwd)
	if err != nil {
		t.Error(err)
	}
	matched, err := encoder.Matches(passwd, passwordHash)
	if err != nil {
		t.Error(err)
	}
	if !matched {
		t.Error("not matched:" + passwordHash)
	}
}
