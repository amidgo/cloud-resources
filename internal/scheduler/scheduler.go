package scheduler

import (
	"context"
	"github.com/amidgo/cloud-resources/internal/executor"
	"log"
	"sync"
	"time"
)

// объект который вызывает метод Execute у каждого элемента из списка executors
// вызов происходит раз в tickTime
type Scheduler struct {
	tickTime  time.Duration
	executors []executor.Executor
}

func New(
	tickTime time.Duration,
	executors ...executor.Executor,
) *Scheduler {
	return &Scheduler{
		tickTime:  tickTime,
		executors: executors,
	}
}

func (s *Scheduler) ExecuteAll(ctx context.Context) {
	// берем контекст с таймаутом в половину tickTime
	ctx, cancel := context.WithTimeout(ctx, s.tickTime/2)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(s.executors))
	for _, exec := range s.executors {
		go func(exec executor.Executor) {
			defer wg.Done()
			err := exec.Execute(ctx)
			log.Printf("execute error %s", err)
		}(exec)
	}
	wg.Wait()
}

// запуск нашего планировщика
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.tickTime)
	// инициализирующий вызов
	s.ExecuteAll(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.ExecuteAll(ctx)
		}
	}
}
