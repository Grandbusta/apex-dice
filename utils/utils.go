package utils

import (
	"math/rand"
	"time"
)

func GenerateTargetNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(13-2) + 2
}

func GenerateDiceNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(7-1) + 1
}
