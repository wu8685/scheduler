package executor

import (
	"testing"
	"time"
)

type TestTask struct {
	logs *[]string
}

func (t TestTask) Do() error {
	*t.logs = append(*t.logs, "executed")
	return nil
}

func (t TestTask) Interrupt() error {
	*t.logs = append(*t.logs, "interrupted")
	return nil
}

func TestRunAtStartPoint(t *testing.T) {
	logs := []string{}
	task := TestTask{&logs}
	startPoint := time.Now().Add(1 * time.Second)

		exe := NewExecutor(startPoint, 2*time.Second)
		exe.Register(task)

	<- time.After(2*time.Second)
	exe.Stop()
	if len(logs) != 2 {
		t.Fatalf("expected hasing 2 log but %d", len(logs))
	}
	
	if logs[0] != "executed" {
		t.Fatalf(`expected hasing log %s but %s`, "executed", logs[0])
	}
	
	if logs[0] != "interrupted" {
		t.Fatalf(`expected hasing log %s but %s`, "interrupted", logs[0])
	}
}
