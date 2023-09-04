package utils

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func Test_viperEnvVariable(t *testing.T) {
	str, err := ViperEnvVariable("DB_LOG")
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}
	//t.Logf(str)
	log.Fatal(str)

}

func Test_RunCommand(t *testing.T) {

	ctx := context.Background()
	str := "pwd"
	err := RunCommand(ctx, str)
	if err != nil {
		t.Error(err)
	}

	//log.Fatal(err)
}
