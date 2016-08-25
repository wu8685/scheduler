package executor

import (
	"fmt"
	"time"
	"reflect"

	"github.com/wu8685/scheduler/task"
)

type Executor struct {
	startPoint time.Time
	interval   time.Duration

	tasks []task.Task
	stop chan time.Duration
}

func NewExecutor(startDate time.Time, interval time.Duration) (*Executor) {
	exe := &Executor{
		startDate,
		interval,
		[]task.Task{},
		make(chan time.Duration, 1),
	}
	runForever(exe.Run)
	return exe
}

func (exe *Executor) Register(t task.Task) *Executor {
	exe.tasks = append(exe.tasks, t)
	return exe
}

func (exe *Executor) Run() error {
	toStart := exe.startPoint.Sub(time.Now())
	if toStart <= 0 {
		return fmt.Errorf("start date should not before now.")
	}

	select {
		case <-time.After(toStart):
		exe.runTasks()
		case timeout:=<-exe.stop:
		exe.interruptTasks(timeout)
		return nil
	}
	
	timer := time.NewTimer(exe.interval)
	for {
		select {
			case <-timer.C:
			exe.runTasks()
			case timeout:=<-exe.stop:
			exe.interruptTasks(timeout)
			return nil
		}
	}
}

func (exe *Executor) Stop(timeout time.Duration) {
	exe.stop<-timeout
}

func (exe *Executor) runTasks() {
	for _, t := range exe.tasks {
		runForever(t.Do)
	}
}

func (exe *Executor) interruptTasks(timeout time.Duration) {
	chans := []reflect.SelectCase{}
	for _, t := range exe.tasks {
		finish := runUnderWatch(t.Interrupt)
		chans = append(chans, reflect.SelectCase{
			Dir: reflect.SelectRecv,
			Chan: reflect.ValueOf(finish),
		})
	}
	
	required := len(chans)
	for {
		select {
			case reflect.Select(chans):
			required--
			if required == 0 {
				return
			}
			case <-time.After(timeout):
			return
		}
	}
}