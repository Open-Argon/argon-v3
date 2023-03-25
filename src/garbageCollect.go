package main

import (
	"runtime"
	"time"
)

func garbageCollect() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			runtime.GC()
		}
	}()
}
