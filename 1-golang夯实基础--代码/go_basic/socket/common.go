package transport

import (
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		log.Printf("error: %s", err)
		os.Exit(1)
	}
}

var MAGIC = []byte{1, 1, 5, 2, 0}

type AddRequest struct {
	RequestId int
	A         int
	B         int
}

type AddResponse struct {
	RequestId int
	Sum       int
}
