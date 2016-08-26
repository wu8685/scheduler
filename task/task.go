package task

type Task interface {
	Do()
	Interrupt()
}
