package pkg

import (
	"strings"
	"github.com/jsimonetti/pwscheme/ssha"
	"github.com/jsimonetti/pwscheme/ssha256"
	"github.com/jsimonetti/pwscheme/ssha512"
)

func Verify(password, hash string) bool {
	if strings.HasPrefix(hash, "{SSHA}") {
		valid, _ := ssha.Validate(password, hash)
		return valid
	}
	if strings.HasPrefix(hash, "{SSHA256}") {
		valid, _ := ssha256.Validate(password, hash)
		return valid
	}
	if strings.HasPrefix(hash, "{SSHA512}") {
		valid, _ := ssha512.Validate(password, hash)
		return valid
	}

	return false
}
