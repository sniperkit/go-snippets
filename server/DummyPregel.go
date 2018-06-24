package main

//import (
//"log"
//"project_c9f7_i5l8_o0p4_p0j8/msg"
//"project_c9f7_i5l8_o0p4_p0j8/pregel/job"
//)

//// this file is for testing purposes

////============================================================
//// Dummy method to test message passing between client - server - workers.
//// This MUST send a return value through cDone before returning!
//// This function is NOT responsible for closing the channels.
//// Pregel is responsible for doing timeout
//// Pregel is responsible for doing partitioning of the given workers.
//// DUMMY FUNCTION

//func DummySuperstep(request Request,
//selectedWorkers []int,
//cIn chan msg.FromWorker,
//cOut chan msg.FromServer,
//cResult chan job.Result) {
//requestId := request.RequestId
//clientId := request.ClientId
//numMessages := 2
//counter := 0
//log.Printf("Superstep() c:%v, r:%v start. %v messages\n", requestId, clientId, numMessages)

//for i := 0; i < numMessages; i++ {
//for _, worker := range selectedWorkers {
//vm := msg.FromServer{msg.V2V}

//vm.Msg = counter
//vm.DstWorker = worker
//vm.StepNum = counter
//cOut <- vm
//}
//counter = counter + 1
//}

//for i := 0; i < numMessages*len(selectedWorkers); i++ {
//msg := <-cIn
//log.Printf("Superstep recevied %v\n", msg)
//}

//log.Printf("\nSuperstep() c:%v, r:%v done.\n", clientId, requestId)
//cResult <- pregel.SUCCESS
//return
//}

//// func DummySuperstep(request Request,
////                selectedWorkers []int,
////                cIn chan pregel.Message,
////                cOut chan pregel.Message,
////                cDone chan job.Result) {
////   requestId := request.RequestId
////   clientId := request.ClientId
////   log.Printf("Superstep() c:%v, r:%v start.\n", requestId, clientId)

////   // Send some messages for a couple work steps.
////   numWorkSteps := 2
////   numMessages := 2
////   for j := 0; j < numWorkSteps; j++ {
////     counter := 0
////     for i := 0; i < numMessages; i++ {
////       for _, worker := range selectedWorkers {
////         var vm pregel.Message
////         vm.FromId = -1
////         vm.Value = counter
////         vm.ToId = worker
////         vm.Superstep = counter
////         // log.Printf("Superstep() c:%v, r:%v sending: %v\n", clientId, requestId, vm)
////         cOut <- vm
////       }
////       counter = counter + 1
////     }

////     // Receive responses from each worker.
////     for i := 0; i < numMessages*numWorkSteps; i++ {
////       // InMsg := <- cIn
////       <- cIn
////       // log.Printf("Superstep() c:%v, r:%v receiving: %v\n", clientId, requestId, InMsg)
////     }
////   }

////   log.Printf("\nSuperstep() c:%v, r:%v done.\n", clientId, requestId)
////   if (request.IterationCount == 2) {
////     cDone <- 1
////   } else {
////     cDone <- 2
////   }
////   return
//// }
