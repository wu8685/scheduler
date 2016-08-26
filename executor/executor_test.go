package executor

import (
	"testing"
	"time"
)

type TestTask struct {
	logs *[]string
}

func (t TestTask) Do() {
	*t.logs = append(*t.logs, "executed")
}

func (t TestTask) Interrupt() {
	*t.logs = append(*t.logs, "interrupted")
}

func TestRunAtStartPoint(t *testing.T) {
	logs := []string{}
	task := TestTask{&logs}
	startPoint := time.Now().Add(1 * time.Second)

	exe := NewExecutor(startPoint, 2*time.Second)
	exe.Register(task)

	<-time.After(2 * time.Second)
	exe.Stop()
	<-time.After(1 * time.Second)
	if len(logs) != 2 {
		t.Fatalf("expected hasing 2 log but %d", len(logs))
	}

	if logs[0] != "executed" {
		t.Fatalf(`expected hasing log %s but %s`, "executed", logs[0])
	}

	if logs[1] != "interrupted" {
		t.Fatalf(`expected hasing log %s but %s`, "interrupted", logs[1])
	}
}

func TestRunAtLoop(t *testing.T) {
	logs := []string{}
	task := TestTask{&logs}
	startPoint := time.Now().Add(1 * time.Second)

	exe := NewExecutor(startPoint, 1*time.Second)
	exe.Register(task)

	<-time.After(3 * time.Second)
	exe.Stop()
	<-time.After(1 * time.Second)
	if len(logs) != 3 {
		t.Fatalf("expected hasing 3 log but %d: %v", len(logs), logs)
	}

	if logs[0] != "executed" {
		t.Fatalf(`expected hasing log %s but %s`, "executed", logs[0])
	}

	if logs[1] != "executed" {
		t.Fatalf(`expected hasing log %s but %s`, "executed", logs[1])
	}

	if logs[2] != "interrupted" {
		t.Fatalf(`expected hasing log %s but %s`, "interrupted", logs[2])
	}
}
