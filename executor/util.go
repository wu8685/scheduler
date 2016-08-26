package executor

func runOnce(runner func()) {
	go func() {
		runner()
	}()
}
