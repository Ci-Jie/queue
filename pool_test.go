package queue

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPoolWithQueueTask(t *testing.T) {
	totalN := 5
	taskN := 100
	rets := make(chan struct{}, taskN)

	p := NewPool(totalN)
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, totalN, p.Workers())

	for i := 0; i < taskN; i++ {
		assert.NoError(t, p.QueueTask(func(context.Context) error {
			rets <- struct{}{}
			return nil
		}))
	}

	for i := 0; i < taskN; i++ {
		<-rets
	}

	// shutdown all, and now running worker is 0
	p.Release()
	assert.Equal(t, 0, p.Workers())
}
