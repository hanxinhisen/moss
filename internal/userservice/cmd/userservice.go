// Created by Hisen at 2022/3/4.
package main

import (
	"github.com/hanxinhisen/moss/internal/userservice"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	userservice.NewApp("cmdb").Run()
}
