package backgroundconsumers

import (
	"context"
	"sync"
)

type BgConsumer interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
}
