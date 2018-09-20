package xrdb

import (
	"log"
	"testing"
)

func TestOption(t *testing.T) {
	log.Println(HMSet("aaa", "haha", "12", "100"))
	log.Println(HMSet("aaa", "yaya", "13", "100"))
	// HDel("aaa", "yaya")
	log.Println(HMGet("aaa", "yaya"))
}
