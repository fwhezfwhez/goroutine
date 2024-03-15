package goroutine

import (
	"encoding/json"
	"time"
)

// gs is wraped object representing a goroutine
type Gs struct {
	f                    func()
	unqKey               string
	expectedExpireSecond int64

	callerStack []string

	startAt time.Time // start_at

	recvFinish chan struct{}

	firstCaller string
}

// function f in gs has done
func (gs *Gs) done() {
	gs.recvFinish <- struct{}{}
}

// returns expire_at and whether will expire
// false, this goroutine will never die
func (gs *Gs) ExpireAt() (time.Time, bool) {
	if gs.expectedExpireSecond == -1 {
		return time.Time{}, false
	}

	return gs.startAt.Add(time.Duration(gs.expectedExpireSecond) * time.Second), true
}

// info of a gs
func (gs *Gs) Info() string {

	expireAt, willExpire := gs.ExpireAt()

	b, _ := json.MarshalIndent(map[string]interface{}{
		"unique_key": gs.unqKey,
		"start_at":   gs.startAt,
		"expires":    map[string]interface{}{"will_expire": willExpire, "expire_at": expireAt},
		"caller":     gs.callerStack,
	}, "  ", "  ")

	return string(b)
}
func (gs *Gs) UniqKey() string {
	return gs.unqKey
}

// Get a read-only view of gs routine
func (gs Gs) GetView() GsView {

	expireAt, willExpire := gs.ExpireAt()
	return GsView{
		UnqKey:      gs.unqKey,
		WillExpire:  willExpire,
		ExpireAt:    expireAt,
		CallerStack: gs.callerStack,
		StartAt:     gs.startAt,
	}
}

// 返回协程的初始发起位置
func (gs Gs) GetLoc() string {
	return gs.firstCaller
}
