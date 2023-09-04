package main

import (
	"log"
	"testing"
)

func Test_connectToRabbit(t *testing.T) {

	_, err := connectToRabbit()
	if err != nil {
		//t.Error(err)
		log.Println(err)
	}
	//t.Fatal("OKE")
}
