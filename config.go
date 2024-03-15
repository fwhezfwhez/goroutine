package goroutine

import (
	"fmt"
	"time"
)

// Event callback
// All callbacks below can be well set by users
var (
	// handle panic event
	HandlePanic func(e interface{}) = func(e interface{}) {
		fmt.Printf("recv panic %v \n", e)
	}

	// handle a bad goroutine
	HandleBadGs func(gs *Gs) = func(gs *Gs) {
		fmt.Printf("goroutine timeout first_caller=%s \n", gs.firstCaller)
	}
)

// Profile arguments
var (
	// How long a zombie info will be storaged in gss.zombies
	ZombieStorageSeconds int = 30

	// In which interval to clear expired zombies
	ZombieClearInterval time.Duration = 30 * time.Second

	// If zombieIndex arrives it, will stop setting zombie gs and forcely trigger clear
	ZombieMaxBustNum int = 5000000
)
