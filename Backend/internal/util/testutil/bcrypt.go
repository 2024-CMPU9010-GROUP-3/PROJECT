package testutil

import (
	"regexp"

	"github.com/pashagolub/pgxmock/v4"
	"golang.org/x/crypto/bcrypt"
)


type bcryptArgument struct {
	expected string
}

var bcryptHashPattern = regexp.MustCompile(`^\$2[ayb]\$.{56}$`)

func (b bcryptArgument) Match(arg interface{}) bool {
	hash, ok := arg.(string)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(b.expected))
	return ok && bcryptHashPattern.MatchString(hash) && err == nil
}

func BcryptArg(expected string) pgxmock.Argument {
	return bcryptArgument{expected: expected}
}