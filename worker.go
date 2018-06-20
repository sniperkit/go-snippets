package main


func smsWorker() error {
	server, err := startServer()
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker("sms_worker", 0)

	return worker.Launch()
}
