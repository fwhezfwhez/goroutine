package goroutine

import (
	"sync"
	"testing"
)

// BenchmarkGs-4   	  300000	      4606 ns/op
func BenchmarkGs(b *testing.B) {
	wg := sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		ProtectedGo(func() {
			defer wg.Done()
			_ = 5
		}, GoParam{
			UnqKey:               "test_benchmark_gs",
			ExpectedExpireSecond: 1,
			ShouldProtected:      true,
		})
	}
	wg.Wait()
}

// BenchmarkGoroutine-4   	 5000000	       263 ns/op
func BenchmarkGoroutine(b *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = 5
		}()
	}

	wg.Wait()
}
