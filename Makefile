build:
	go build

worker:
	./machinery-message-center-demo sms-worker

send:
	./machinery-message-center-demo send-sms