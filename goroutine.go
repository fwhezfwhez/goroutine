package goroutine

import (
	"fmt"
	"math"
	"runtime"
	"sync/atomic"
	"time"
)

// ProtectedGo param
type GoParam struct {
	// A goroutine will mark this unique key and saved.
	// Will auto-generated if emtpy.
	UnqKey string

	// How long you expect this goroutine to die.
	// -1, this goroutine should never die.
	// >0, this goroutine should die in [ExpectedExpireSecond] seconds interval
	ExpectedExpireSecond int64

	// caller stack depth. default 1.
	captureLen int
	// where caller is
	callerStack []string

	firstCaller string

	// If shouldProtected is false, then it's just a simply goroutine
	// If true, gss will be add to gss-management, time-out err and zombie monitor
	ShouldProtected bool
}

// Go with exist goroutine
func (gp GoParam) begin(start string) GoParam {

	if gp.firstCaller == "" {
		gp.firstCaller = start
	}

	if gp.ShouldProtected == false {
		return gp
	}

	if len(gp.callerStack) == 0 {
		gp.callerStack = []string{
			start,
		}
	} else {
		gp.callerStack = append([]string{start}, gp.callerStack...)
	}

	return gp
}

// help generate uniqKey of a goroutine
var offset int32

// Start a goroutine in protected mod
// This function should be straightly called  where you put `go` and should not be wrap by functions again,because if wrapped, caller stack might wrong depth.
// If consider wrapped, call Go(f ,param...) instead.
// If consider wrapped in many layer , call GoDepth(depth, f, param...)
func ProtectedGo(f func(), params ...GoParam) {
	here := SpotHere()
	if len(params) >= 1 {
		go protectedGo(f, params[0].begin(here))
	} else {
		go protectedGo(f, GoParam{}.begin(here))
	}
}

// If want to use ProtectedGo wrapped in your pkg, you should call Go or Godepth to help locate call stack correctly.
func Go(f func(), params ...GoParam) {
	here := SpotHereV2(3)
	if len(params) >= 1 {
		go protectedGo(f, params[0].begin(here))
	} else {
		go protectedGo(f, GoParam{}.begin(here))
	}
}

func GoDepth(depth int, f func(), params ...GoParam) {
	here := SpotHereV2(depth)
	if len(params) >= 1 {
		go protectedGo(f, params[0].begin(here))
	} else {
		go protectedGo(f, GoParam{}.begin(here))
	}
}

func protectedGo(f func(), param GoParam) {
	if offset >= math.MaxInt32-1000 {
		offset = 0
	}

	// handle panic in goroutine
	defer func() {
		if e := recover(); e != nil {
			if HandlePanic != nil {
				HandlePanic(fmt.Errorf("panic from %v loc=%s", e, param.firstCaller))
			}
			// fmt.Printf("panic recover from %v \n %s", e, debug.Stack())
		}
	}()

	if param.ShouldProtected == false {
		f()
		return
	}

	// handle goroutine key
	if param.UnqKey == "" {
		param.UnqKey = fmt.Sprintf("%d:%d", time.Now().UnixNano(), atomic.AddInt32(&offset, 1))
	}

	// handle stack
	if param.captureLen == 0 {
		param.captureLen = 1
	}
	if param.callerStack == nil {
		param.callerStack = make([]string, 0, 10)
	}

	for i := 0; i < param.captureLen; i++ {
		_, file, l, _ := runtime.Caller(i)
		stackline := fmt.Sprintf("%s:%d", file, l)
		param.callerStack = append(param.callerStack, stackline)
	}

	// gs is the wrapped goroutine
	gs := &Gs{
		f:                    f,
		unqKey:               param.UnqKey,
		expectedExpireSecond: param.ExpectedExpireSecond,
		callerStack:          param.callerStack,
		startAt:              time.Now(),
		recvFinish:           make(chan struct{}, 1),
		firstCaller:          param.firstCaller,
	}
	_ = gs

	func(gs *Gs) {
		gsSchedules.addGs(gs)

		defer gs.done()
		gs.f()
	}(gs)
}
