package main

import (
	"fmt"
	"os"
	"sync"
)

var debug = os.Getenv("__ARGON_DEBUG__") == "true"

var __debugPrints = [][]any{}
var __debugPrintsLock = sync.RWMutex{}

func debugInit() {
	if debug {
		fmt.Println("In debug mode...")
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("debugPrintln: panic:", r)
				}
			}()
			for {
				__debugPrintsLock.RLock()
				for _, v := range __debugPrints {
					fmt.Println(v...)
				}
				__debugPrintsLock.RUnlock()
				__debugPrintsLock.Lock()
				__debugPrints = [][]any{}
				__debugPrintsLock.Unlock()
			}
		}()
	}
}

func debugPrintln(a ...any) {
	if debug {
		__debugPrintsLock.Lock()
		__debugPrints = append(__debugPrints, a)
		__debugPrintsLock.Unlock()
	}
}
