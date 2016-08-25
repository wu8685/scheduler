package executor

func runForever(runner func() error) {
	go func() {
		runner()
	} ()
}

func runUnderWatch(runner func() error) chan struct{} {
	finish := make(chan struct{})
	go func() {
		runner()
		finish<-struct{}{}
	} ()
	return finish
} 