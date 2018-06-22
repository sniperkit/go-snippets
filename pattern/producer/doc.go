// Package producer describes Producer/Consumer pattern
// Producer generates data and pushes to producer channel
// and waits consumer to get, consumer get data from producer channel
// and feedback to producer. There needs two channel, one for data, one for feedback
// ref, http://www.golangpatterns.info/concurrency/producer-consumer
package producer
