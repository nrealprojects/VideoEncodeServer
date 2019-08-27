package main

import (
	"fmt"
	"projects/VideoEncodeServer/models"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	ch := make(chan models.PackData, 100)
	go func() {
		for i := 0; i < 10; i++ {
			pack := Get()
			pack.Init(1300)
			ch <- pack
		}
	}()

	go func() {
		for {
			select {
			case pack := <-ch:
				fmt.Println(pack.ToString())
				Put(pack)
			case <-time.After(1 * time.Second):
				return
			}
		}
	}()

	time.Sleep(2 * time.Second)
}
