package htpasswd

import "golang.org/x/crypto/bcrypt"

// Bcrypt generates a hash and returns it.
func Bcrypt(pw []byte, cost int) (hash []byte, err error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	hash, err = bcrypt.GenerateFromPassword(pw, cost)
	if err != nil {
		return hash, err
	}
	return hash, nil
}
