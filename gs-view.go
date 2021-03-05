package goroutine

import "time"

type GsView struct {
	UnqKey     string    `json:"unq_key"`
	WillExpire bool      `json:"will_expire"`
	ExpireAt   time.Time `json:"expire_at"`

	CallerStack []string  `json:"caller_stack"`
	StartAt     time.Time `json:"start_at"`
}
