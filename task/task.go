package task

type Task interface {
	Do() error
	Interrupt() error
}
