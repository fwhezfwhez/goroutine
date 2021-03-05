package goroutine

// Event callback
// All callbacks below can be well set by users
var (
	// handle panic event
	HandlePanic func(e interface{})

	// handle a bad goroutine
	HandleBadGs func(gs *Gs)
)

// Profile arguments
var (
	// How long a zombie info will be storaged in gss.zombies
	// by default 3days
	ZombieStorageSeconds int
)
