package main

import (
	"math/rand"
        "fmt"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!"
func RandStringBytesRmndr(n int) string {
	b := make([]byte,n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main () {

	for i:= 0;i < 25; i++ {
		fmt.Printf("Random string %s\n",RandStringBytesRmndr(14))
	}
}
