package data

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_initDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Error(err)
		//log.Fatal(err)
	}
}

func Test_InsertLog(t *testing.T) {
	ctx := context.Background()
	value := fmt.Sprintf("test4_%s", time.Now().GoString())
	err := InsertLog(ctx, value)
	if err != nil {
		t.Error(err)
		log.Fatal(err)
	}

}

func Test_Insert(t *testing.T) {
	//	ctx := context.Background()
	//	err := InsertLog(ctx, value)
	//
	//	entry := LogEntry{
	//		Name:      "name",
	//		Data:      "data",
	//		CreatedAt: time.Now(),
	//		UpdatedAt: time.Now(),
	//	}
	//
	// err = entry.Insert(entry)
	//
	//	if err != nil {
	//		log.Println(err)
	//		//	t.Error(err)
	//	}
	//
	// log.Println("end test")
}
