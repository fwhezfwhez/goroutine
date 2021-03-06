package goroutine

import (
	"fmt"
	"github.com/fwhezfwhez/cmap"
	"runtime/debug"
	"sync/atomic"
	"time"
)

var gsSchedules = NewGsSchedule()

// gs manager
type GsSchedule struct {
	// will storage each gs routines.
	// When goroutine expires, this gs routien will remove from gsm and be set to zombies
	gsm *cmap.Map

	// When goroutine expires, this gs routien will remove from gsm and be set to zombies
	zombies *cmap.Map

	// When goroutine starts with no expire setting, will be put into eternals.
	// If a goroutine marked eternal however stop quickly, it will also be correctly removed both from gsm and eternals.
	eternals *cmap.Map

	// each time meets a zombie will incr one times
	zombieIndex    int64
	startedZombied bool
}

// new a gs routine manager
func NewGsSchedule() *GsSchedule {
	return &GsSchedule{
		gsm:      cmap.NewMap(),
		zombies:  cmap.NewMap(),
		eternals: cmap.NewMap(),
	}
}

// add a gs routine into gss manager
func (gss *GsSchedule) addGs(gs *Gs) {
	if gss.startedZombied == false {
		gss.zombieD()
	}

	expireAt, willExipire := gs.ExpireAt()

	gss.gsm.Set(gs.unqKey, gs)

	if willExipire {
		go func() {

			var isTimeout bool

		L:
			for {
				select {
				case <-time.After(expireAt.Sub(time.Now())):
					// when gs reaches its deadline, it haven't recv finish yet.
					if HandleBadGs != nil {
						HandleBadGs(gs)
					}
					fmt.Printf("recv a zombie goroutine: %s\n", gs.callerStack)

					// remove it from gss.gsm to gss.zombie
					gss.gsm.Delete(gs.unqKey)
					gss.zombies.SetEx(gs.unqKey, gs, DefaultInt(ZombieStorageSeconds, 5*60))
					atomic.AddInt64(&gss.zombieIndex, 1)

					isTimeout = true
					break L
				case <-gs.recvFinish:
					// gs recv finish signal
					gss.finishGs(gs)
					return
				}
			}
			_ = isTimeout
		}()
		return
	}

	gss.eternals.Set(gs.unqKey, gs)
}

// start a goroutine to keep gss.zombie cleared each intervals.
func (gss *GsSchedule) zombieD() {

	if !gss.startedZombied {
		gss.startedZombied = true
		go func() {
			if e := recover(); e != nil {
				fmt.Printf("panic recover from %v\n %s\n", e, debug.Stack())
			}

			for {
				n := gss.zombies.ClearExpireKeys()
				if n > 0 {
					atomic.AddInt64(&gss.zombieIndex, -int64(n))
				}
				fmt.Printf("clear %d expired zombie keys \n", n)
				if ZombieClearInterval == 0 {
					time.Sleep(30 * time.Minute)
				} else {
					time.Sleep(ZombieClearInterval)
				}

			}
		}()
	}

}

// gss will remove gs if this gs routine is finished
func (gss *GsSchedule) finishGs(gs *Gs) {
	gss.gsm.Delete(gs.unqKey)
	gss.eternals.Delete(gs.unqKey)
}

// Will monitor zombie goroutines with specific size.
func (gss *GsSchedule) Monitoring(size int) MonitorResult {

	var rs = MonitorResult{
		ZombieViews: make([]GsView, 0, 10),
	}

	rs.ZombieCount = gss.zombies.RealLength()

	var offset int
	gss.zombies.Range(func(key string, value interface{}) bool {
		gs := value.(*Gs)
		rs.ZombieViews = append(rs.ZombieViews, gs.GetView())

		offset ++
		if offset >= size {
			return false
		}

		return true
	})
	return rs
}

func (gss *GsSchedule) MonitorOne(unqKey string) MonitorResult {
	var rs = MonitorResult{
		ZombieCount: 0,
		ZombieViews: make([]GsView, 0, 10),
	}

	gsi, exist := gss.zombies.Get(unqKey)

	if !exist {
		return rs
	}

	gs := gsi.(*Gs)

	rs.ZombieViews = append(rs.ZombieViews, gs.GetView())
	rs.ZombieCount = 1
	return rs
}
