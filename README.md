## goroutine
[goroutine] is a pkg which provides safe usage of using golang goroutine.

To get well divided, `goroutine` is described  as golang official object, and `[goroutine]` is described for this repo.

## Declaration
[goroutine] is not redesigning another goroutine. It's just a wrap to make goroutine perform better in group team.

## Why do I design goroutine, rather than using `go func(){}()`?
Golang basically supports goroutine and it's easy to use , like `go handler()`.

It's just easy and convenient, however!

Most of you should have been faced with these kind of questions:

- A goroutine function panics and trigger your application node hung up or restart. [This is extremely fatal!!]
- A goroutine once started, it's hardly being monitored. [There're too many goroutines in pprof/goroutine, hardly locate what you want!]
- Most memory leak (by goroutine bursting) is caused by mis-understanding lifetime of a goroutine, this package force user to estimate how long your goroutine should have. Err will be reported if unexpected alive goroutine occurs.


And thus, using [goroutine] will improve all advantages below:

- Auto recovers if goroutine meets an inner panic.Besides, how to handle a panic can be set by `goroutine.HandlePanic = func(e interface{}) {}`.
- Monitor zombie goroutines.
- Alert and report those zombie goroutines.

Know that using [goroutine] will have some extra cost for monitoring a goroutine, they're:

- An extra goroutine to wait core goroutine to finish or timeout warning.
- Each goroutine will be managed through a gss center, it's a concurrent map (It means two irrelevant goroutines are put into a race container, there should be lock cost too).


In which case you should use [goroutine]:
- Most of your teammates are new. They're ones easier to produce panic.
- Some of your jobs are important and just put into product use (it might cause unexpected leak or panic and thus should be monitored for some time).


## Start
`go get github.com/fwhezfwhez/goroutine`

## Usage
```go
package main

import (
	"fmt"
	"goroutine"
	"math/rand"
	"runtime/debug"
	"time"
)

func init() {
	// config zombie storage time
	goroutine.ZombieStorageSeconds = 3 * 24 * 60 * 60

	// specific method to handle panic in goroutine
	goroutine.HandlePanic = func(e interface{}) {
		fmt.Printf("recv a panic from protected-go: %s \n %s \n", fmt.Sprintf("%v", e), debug.Stack())
	}

	// specific method to alert a bad goroutine.
	goroutine.HandleBadGs = func(gs *goroutine.Gs) {
		fmt.Printf("recv a bad goroutine: %s\n", gs.Info())
	}

	// Each 5 seconds spying on zombie goroutines
	go func() {
		for {
			info := goroutine.GlobalGss().Monitoring(10)
			fmt.Println(info)

			time.Sleep(5 * time.Second)
		}
	}()

}

func main() {

	// normal goroutines with short deadline
	for i := 0; i < 10000; i ++ {
		goroutine.ProtectedGo(func() {
			fmt.Println("normal goroutine")
		}, goroutine.GoParam{
			UnqKey:               fmt.Sprintf("test_protected_go_normal_goroutine_%d", i),
			ExpectedExpireSecond: 1,
			ShouldProtected:      true,
		})
	}

	// An eternal goroutine without deadline
	goroutine.ProtectedGo(func() {
		fmt.Println("an eternal goroutine")
	}, goroutine.GoParam{
		UnqKey:               "test_protected_go_eternal_goroutine",
		ExpectedExpireSecond: -1,
		ShouldProtected:      true,
	})

	// zombie goroutines. It expects 3s to run however lasting for 5 seconds
	for i := 0; i < 1000; i ++ {
		go func(i int) {
			time.Sleep(time.Duration(rand.Intn(10)%15) * time.Second)
			goroutine.ProtectedGo(func() {
				fmt.Println("a being zombie goroutine")
				time.Sleep(5 * time.Second)
			}, goroutine.GoParam{
				UnqKey:               fmt.Sprintf("test_protected_go_zombie_goroutine_%d", i),
				ExpectedExpireSecond: 3,
				ShouldProtected:      true,
			})
		}(i)
	}

	select {}
}

```



