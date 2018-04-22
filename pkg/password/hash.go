package password

import "gopkg.in/hlandau/passlib.v1"

func Hash(password string) (hash string, err error) {
	return passlib.Hash(password)
}
