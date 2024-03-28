package cleaner

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TodoCallback = func(ctx context.Context) error
type Cleaner interface {
	Add(f ...TodoCallback)
	CloseAll(timeout time.Duration)
}

func New(cap uint) Cleaner {
	return &instance{
		mx:     sync.Mutex{},
		todo:   make([]TodoCallback, 0, cap),
		amount: cap,
	}
}

// gracefully shutdowns app
func DefferedClear(c Cleaner, timeout time.Duration) {
	if err := recover(); err != nil {
		fmt.Println(fmt.Sprintf(recoverMessageTemplate, err))
	}

	c.CloseAll(timeout)
}
