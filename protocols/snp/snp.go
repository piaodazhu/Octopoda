package snp

import (
	"math/rand"
	"sync"
	"time"
)

type SerialNumberProcessor struct {
	cap   int
	idx   int
	ringQ []uint32
	set   map[uint32]struct{}
	lock  sync.Mutex
}

var snp SerialNumberProcessor

func init() {
	snp = SerialNumberProcessor{
		cap:   1024,
		idx:   0,
		ringQ: make([]uint32, 1024),
		set:   map[uint32]struct{}{},
	}
	rand.Seed(time.Now().UnixNano())
}

func GenSerial() uint32 {
	serial := rand.Uint32()
	for serial == 0 {
		serial = rand.Uint32()
	}
	return serial
}

func CheckSerial(serial uint32) bool {
	snp.lock.Lock()
	defer snp.lock.Unlock()
	if _, found := snp.set[serial]; found {
		return false 
	}
	snp.set[serial] = struct{}{}
	snp.ringQ[snp.idx] = serial
	snp.idx = (snp.idx + 1) % snp.cap
	if snp.ringQ[snp.idx] != 0 {
		delete(snp.set, snp.ringQ[snp.idx])
	}
	return true 
}
