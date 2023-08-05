package async

import (
	"context"
	"fmt"
	"time"

	"github.com/tkudlicka/portflux-api/app/async/healthchecker"
	"github.com/tkudlicka/portflux-api/config"
)

type async struct {
	config config.Config
}

func New(cfg config.Config) async {
	return async{
		config: cfg,
	}
}

func (a async) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		go healthchecker.Run(ctx, cancel, fmt.Sprintf("http://:%d/health", a.config.Port), a.config.Async.Interval.Duration)

		for ctx.Err() == nil {
			<-time.After(1 * time.Second)
		}
		return fmt.Errorf("async process stopped")
	}
}
