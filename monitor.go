package goroutine

// monitor result.
// expected call monitoring  per 10s
type MonitorResult struct {
	ZombieCount int
	ZombieViews []GsView
}

func GlobalGss() *GsSchedule {
	return gsSchedules
}
