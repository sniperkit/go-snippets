RPC (Remote Procedure Call) is like calling a function, only that it is present remotely on a different server as a service. A service exposes many such functions/procedure to its client. And client requires some way to know what are the functions/procedures exposed by this service and what are their parameters.

This is where Apache Thrift comes in. It has its own "Interface Definition Language" (IDL). In this language you define what are the functions and what are their parameters. And then use Thrift compiler to generate corresponding code for any language of your choice. What this means, is that you can implement a function in java, host it on a server and then remotely call it from python.

Important work a framework like Thrift does is this -
Provide a language agnostic Interface Definition Language
A Compiler to compile this IDL to produce client and server code (in same or separate language as required)
    Compiler generated client code exposes a stub interfaces for these functions. The stub code converts the parameters passed to the function into a binary (serialized) format that can be transported on wire over network. This process is called marshaling. The generated client code never has the actual implementation of the function, hence its called a stub.
    At server, the developer use the Compiler generated server code, to actually implement these functions (i.e. write the actual functionality of the function). Generated server side code receives the binary encoded message from client, converts them back to the corresponding language objects and passes it to the developer implemented function. This is called as unmarshaling. In java for example the Compiler generated server code would be Interface that the developer will implement and also various other classes.
    Similarly the result of a function is converted to binary and send to client.

    For parameters to the function, IDL defines its own set of data structure types like List, Map, Struct or Classes apart from native types like Int, String, Boolean, etc. These are then mapped to corresponding language implementations.

    Thrift is similar to SOAP and CORBA. Since they both are used for RPC and provide their own IDL. CORBA and SOAP generally also has a service discovery broker as a middleware for exposing functions/methods to client. For thrift, we normally use Zookeeper for service discovery.

    REST is different, because it does not have IDL, And uses HTTP Methods like GET, PUT and url patterns to call a remote function and pass parameters. Using HTTP methods and url semantics makes it also language agnostic.

    Messaging queue is entirely different. Because it is mostly used in Publish/Subscribe model. Whereas RPC is Client/Server model.

    In Publish/Subscribe, multiple publishers sends/add a serialized message on a queue. The message format is defined by the publisher and has complete control of it. Their definition is semantically associated to the queue on which they are published, but there is no strict checking for their structure. Subscriber, then knowing the kind of the message a queue will have, subscribes for those messages. Publishers dont know who are the client, and Subscribers dont know who are the producers of the message. They only know what kind of message to publish or consume respectively from a queue. The Publisher and Subscriber is responsible for knowing the right serializer and deserializer.
    This is different in Client/Server RPC, since Client knows (in strict sense) what to pass and Server defines it. And also whom to pass.

    Other library similar to Thrift is Protobuf and Avro.
