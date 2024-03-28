package cleaner

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const recoverMessageTemplate = "\033[1m\033[31mterminating service due to main panic:\033[0m\n - %v"

type instance struct {
	mx     sync.Mutex
	todo   []TodoCallback
	amount uint
}

func (self *instance) Add(f ...TodoCallback) {
	self.mx.Lock()
	self.todo = append(self.todo, f...)
	self.mx.Unlock()
}

func (self *instance) CloseAll(timeout time.Duration) {
	self.mx.Lock()

	whisper("closing/clearing resources")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	defer func() {
		self.todo = make([]TodoCallback, 0, self.amount)
		self.mx.Unlock()
	}()

	var (
		errmsg = make([]error, 0, len(self.todo))
		done   = make(chan bool, 1)
	)

	go func() {
		for _, f := range self.todo {
			whisper("clearing " + functionName(f))
			if err := f(ctx); err != nil {
				errmsg = append(errmsg, err)
			}
		}

		done <- true
	}()

	select {
	case <-done:
		break
	case <-ctx.Done():
		whisper(fmt.Sprintf("shutdown cancelled by context: %v", ctx.Err()))
	}

	if len(errmsg) > 0 {
		whisper(fmt.Sprintf(
			"shutdown done with error(s):\n%s",
			strings.Join(func() []string {

				buffer := make([]string, len(errmsg))
				for i, e := range errmsg {
					buffer[i] = e.Error() + "\n"
				}
				return buffer

			}(), ", "),
		))
	}
}
