This package provide handling exit signals (SIGKILL, SIGTERM, SIGQUIT and Interrupt)

It's not difficult to write, just dead simple to use)  

## Installation
```go get -u gopkg.in/cheggaaa/go-exit.v1```

#### Simple example

```go
package main

import (
	"fmt"
	
	"gopkg.in/cheggaaa/go-exit.v1"
)

func main() {
	// start app in other goroutine
	go myAppStart()

	// wait stop signal
	sig := exit.Wait()
	fmt.Printf("'%v' recieved\n", sig)

	// and do app exit
	myAppStop()
}

// start your app here
func myAppStart() {
	fmt.Println("my app start")
	fmt.Println("send 'kill pid' or press Ctrl+C for exit")
}

func myAppStop() {
	// Do something for exit 
	fmt.Println("my app stop")
}

```

#### Advanced example

```go
package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
	
	"gopkg.in/cheggaaa/go-exit.v1"
)

func main() {
	// we can enable http profiling ("net/http/pprof")
	if e := exit.EnableHttpProfiling(""); e != nil {
		fmt.Println("Oops can't enable pprof:", e)
	}
	
	// register callbacks (optionally)
	exit.On(myCallback)

	// start app in other goroutine
	go myAppStart()

	// wait stop signal
	sig := exit.Wait()
	fmt.Printf("'%v' recieved\n", sig)

	// we can determine reason for exit
	switch sig.(type) {
	case os.Signal:
		switch sig.(os.Signal) {
		case os.Interrupt:
			fmt.Println("'ctrl+c' or 'kill -2'")
		case syscall.SIGTERM:
			fmt.Println("maybe kill")
		default:
			fmt.Printf("os.Signal: '%v'\n", sig)
		}
	case string:
		fmt.Println("exit.Exit() was called")
	}

	// and do app exit
	myAppStop()
}

func myAppStart() {
	fmt.Println("my app start")
	fmt.Println("send 'kill pid' or press Ctrl+C for exit")
	// wait 30 seconds and do exit
	time.Sleep(time.Second * 30)
	exit.Exit("Work done!")
}

func myCallback() {
	fmt.Println("my callback")
}

func myAppStop() {
	// Do something for exit 
	fmt.Println("my app stop")
}

```
