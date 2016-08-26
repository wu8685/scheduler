package executor

import (
	"fmt"
	"time"

	"github.com/wu8685/scheduler/task"
)

type Executor struct {
	startPoint time.Time
	interval   time.Duration

	tasks []task.Task
	stop  chan struct{}
}

func NewExecutor(startDate time.Time, interval time.Duration) *Executor {
	exe := &Executor{
		startDate,
		interval,
		[]task.Task{},
		make(chan struct{}, 1),
	}
	runOnce(exe.Run)
	return exe
}

func (exe *Executor) Register(t task.Task) *Executor {
	exe.tasks = append(exe.tasks, t)
	return exe
}

func (exe *Executor) Run() {
	toStart := exe.startPoint.Sub(time.Now())
	if toStart <= 0 {
		fmt.Println("start date should not before now.")
		return
	}

	select {
	case <-time.After(toStart):
		exe.runTasks()
	case <-exe.stop:
		exe.interruptTasks()
		return
	}

	times := 0
	for {
		times++
		duration := exe.startPoint.Add(time.Duration(times) * exe.interval).Sub(time.Now())
		if duration < 0 {
			// If passed the timing, just ignore this time
			continue
		}
		select {
		case <-time.After(duration):
			exe.runTasks()
		case <-exe.stop:
			exe.interruptTasks()
			return
		}
	}
}

func (exe *Executor) Stop() {
	exe.stop <- struct{}{}
}

func (exe *Executor) runTasks() {
	for _, t := range exe.tasks {
		runOnce(t.Do)
	}
}

func (exe *Executor) interruptTasks() {
	for _, t := range exe.tasks {
		runOnce(t.Interrupt)
	}
}
