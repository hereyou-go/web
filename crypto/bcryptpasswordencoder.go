package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordEncoder struct{}

func (encoder *bcryptPasswordEncoder) Gensalt() (salt string, err error) {
	rnd := make([]byte, 16)
	_, err = rand.Read(rnd)
	if err != nil {
		return "", err
	}
	hashed, err := bcrypt.GenerateFromPassword(rnd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (encoder *bcryptPasswordEncoder) Encode(rawPassword string, salt ...string) (string, error) {
	//salt is ignored

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (encoder *bcryptPasswordEncoder) Matches(rawPassword, encodedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(encodedPassword), []byte(rawPassword))

	return err == nil, err
}

var BCrypt = &bcryptPasswordEncoder{}

func init() {
	RegisterPasswordEncoder("bcrypt", BCrypt)
}
