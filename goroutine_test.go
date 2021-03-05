package goroutine

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGoRoutine(t *testing.T) {

	HandleBadGs = func(gs *Gs) {
		fmt.Printf("routine:%s bad lifetime:\n %s \n", gs.UniqKey(), gs.Info())
	}

	ProtectedGo(func() {
		time.Sleep(6 * time.Second)
	}, GoParam{
		ExpectedExpireSecond: -1,
	})

	time.Sleep(15 * time.Second)
}

func TestMonitor(t *testing.T) {

	HandleBadGs = func(gs *Gs) {
		fmt.Printf("routine:%s bad lifetime:\n %s \n", gs.UniqKey(), gs.Info())
	}

	ProtectedGo(func() {
		time.Sleep(50 * time.Second)
	}, GoParam{
		ExpectedExpireSecond: 12,
	})
	ProtectedGo(func() {
		time.Sleep(50 * time.Second)
	}, GoParam{
		ExpectedExpireSecond: 15,
	})
	ProtectedGo(func() {
		time.Sleep(50 * time.Second)
	}, GoParam{
		ExpectedExpireSecond: 10,
	})

	time.Sleep(1 * time.Second)
	time.Sleep(29 * time.Second)


	rs := gsSchedules.Monitoring(4)

	fmt.Println(JSON(rs))

}

func TestRWMutext(t *testing.T) {
	l := sync.RWMutex{}

	go func() {
		time.Sleep(4 * time.Second)
		l.Lock()
		fmt.Println(11111)
		l.Unlock()
	}()

	go func() {
		l.RLock()
		time.Sleep(15 * time.Second)
		l.RUnlock()
	}()

	time.Sleep(50 * time.Second)
}
