package security

import (
	"fmt"
	"time"
)

type Token struct {
	Raw    []byte
	Serial int64
	Age    int
}

var TokenEnabled bool

var tokens [2]Token
var lastfetch time.Time

const Fetchinterval time.Duration = time.Minute

func UpdateTokens(cur, prev Token) bool {
	lastfetch = time.Now()
	tokens = [2]Token{cur, prev}
	return true
}

func chooseToken() (Token, error) {
	if time.Since(lastfetch) > Fetchinterval*2 {
		return Token{}, fmt.Errorf("tokens are not fetched in time")
	}
	cur, prev := tokens[0], tokens[1]
	// current token haven't been fetched.
	if cur.Serial == 0 {
		return Token{}, fmt.Errorf("current token haven't been fetched")
	}

	// current token is stable, choose it.
	if cur.Age > 2 * int(Fetchinterval.Seconds()) {
		return cur, nil
	}

	// previous token doesn't exist, return the current
	if prev.Serial == 0 {
		return cur, nil
	} else {
		return prev, nil
	}
}

func matchToken(serial int64) (Token, error) {
	if tokens[0].Serial == serial {
		return tokens[0], nil
	} else if tokens[1].Serial == serial {
		return tokens[1], nil
	} else {
		return Token{}, fmt.Errorf("no matching token: %d", serial)
	}
}
