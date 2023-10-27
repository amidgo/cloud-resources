package scheduler_test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/amidgo/cloud-resources/internal/executor"
	"github.com/amidgo/cloud-resources/internal/scheduler"
	"sync/atomic"
	"testing"
	"time"
)

type FakeExecutor struct {
	t            *testing.T
	executeDelay time.Duration
	executeCount atomic.Int32
	err          error
}

func NewFakeExecutor(t *testing.T, delay time.Duration, count int, err error) executor.Executor {
	executor := &FakeExecutor{
		executeDelay: delay,
		executeCount: atomic.Int32{},
		err:          err,
	}
	executor.executeCount.Add(int32(count))
	t.Cleanup(func() {
		if executor.executeCount.Load() > 0 {
			t.Fatalf("no calls")
		}
	})
	return executor
}

func (f *FakeExecutor) Execute(ctx context.Context) error {
	if f.executeCount.Load() == 0 {
		f.t.Fatal("unexpected execute")
	}
	select {
	case <-ctx.Done():
		return errors.New("ctx done")
	case <-time.Tick(f.executeDelay):
		f.executeCount.Add(-1)
		return f.err
	}
}

func Test_Scheduler(t *testing.T) {
	cases := []struct {
		executorsCount int
		executeData    struct {
			delay time.Duration
			count int
			err   error
		}
		tickTime time.Duration
	}{
		{
			executorsCount: 3,
			executeData: struct {
				delay time.Duration
				count int
				err   error
			}{
				delay: time.Millisecond * 100,
				count: 5,
				err:   sql.ErrNoRows,
			},
			tickTime: time.Second,
		},
	}

	for _, cs := range cases {
		executors := make([]executor.Executor, 0, cs.executorsCount)
		for i := 0; i < cs.executorsCount; i++ {
			executors = append(executors, NewFakeExecutor(t, cs.executeData.delay, cs.executeData.count, cs.executeData.err))
		}
		scheduler := scheduler.New(cs.tickTime, executors...)
		ctx, cancel := context.WithTimeout(context.Background(), cs.tickTime*time.Duration(cs.executeData.count-1)+(cs.tickTime/2))
		scheduler.Run(ctx)
		cancel()
	}
}
