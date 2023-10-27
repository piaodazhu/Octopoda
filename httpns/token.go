package main

import (
	"math/rand"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			sb.WriteByte(letters[idx])
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return sb.String()
}

type Token struct {
	raw    string
	birth  time.Time
	serial int64
}

var token = [2]Token{
	{
		raw:    "none",
		birth:  time.Now(),
		serial: 0,
	},
	{
		raw:    "none",
		birth:  time.Now(),
		serial: 0,
	},
}

func UpdateTokens() {
	newToken := Token{
		raw:    randStr(16),
		birth:  time.Now(),
		serial: rand.Int63(),
	}
	token = [2]Token{newToken, token[0]}
}

var rollingInterval = time.Minute * 30

func startRollingToken() {
	ticker := time.NewTicker(rollingInterval)
	go func() {
		UpdateTokens()
		for range ticker.C {
			UpdateTokens()
		}
	}()
}
