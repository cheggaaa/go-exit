package exit

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestOnExit(t *testing.T) {
	var aWorks int32
	var a = func() {
		atomic.AddInt32(&aWorks, 1)
	}
	var bWorks int32
	var b = func() {
		atomic.AddInt32(&bWorks, 1)
	}

	On(a)
	go Wait()
	On(b)
	Exit("test")
	if atomic.LoadInt32(&aWorks) != 1 || atomic.LoadInt32(&bWorks) != 1 {
		t.Error("Callbacks do not works")
	}
}

func TestExit(t *testing.T) {
	var sig = make(chan interface{})
	go func() {
		sig <- Wait()
	}()
	Exit("unique test message")
	select {
	case <-time.After(time.Millisecond):
		t.Errorf("Signal timeout")
	case s := <-sig:
		if _, ok := s.(string); !ok || s.(string) != "unique test message" {
			t.Errorf("Unexpected signal: %s", s)
		}
	}
}

func TestEnableHttpProfiling(t *testing.T) {
	// test error
	var e error
	if e = EnableHttpProfiling(":-80"); e == nil {
		t.Errorf("Expected error, but recieve nil")
	}
	t.Log(e)
}
