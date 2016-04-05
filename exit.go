// Small and simple helper for handling exit signals (SIGKILL, SIGTERM, SIGQUIT and Interrupt)
package exit

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var closeSignal = make(chan string)
var onExitCallbacks = make([]func(), 0)
var onExitMu sync.Mutex

// Just wait the signals for exit
// Return catched signal
func Wait() interface{} {
	// catch signals
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	var sig interface{}
	select {
	case sig = <-exit:
	case sig = <-closeSignal:
	}
	onExitCall()
	return sig
}

// Send exit signal with given message
func Exit(message string) {
	closeSignal <- message
}

// Helper for enable http profiling
// Default addr is ':6060'
func EnableHttpProfiling(addr string) (err error) {
	if addr == "" {
		addr = ":6060"
	}
	var e = make(chan error)
	go func() {
		e <- http.ListenAndServe(addr, nil)
	}()
	select {
	case err = <-e:
		return
	case <-time.After(time.Millisecond):
		return nil
	}
	return
}

// Add callback for exit
func On(f func()) {
	onExitMu.Lock()
	onExitCallbacks = append(onExitCallbacks, f)
	onExitMu.Unlock()
}

func onExitCall() {
	onExitMu.Lock()
	defer onExitMu.Unlock()
	for _, f := range onExitCallbacks {
		f()
	}
}
