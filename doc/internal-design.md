# Notes on qiloop implementation

## EndPoint abstraction

The package bus/net offers an abstraction of an established network connection
called EndPoint. The EndPoint interface allow multiple parties to send and
receive messages.

	type EndPoint interface {
		Send(m Message) error
		AddHandler(f Filter, c Consumer, cl Closer) int
		RemoveHandler(id int) error
		[...]
	}


Sending a message is done using the synchronous method `Send`. In order to
receive a message one needs to register an handler composed of 3 methods:

- a `Filter` method to select the messages.
- a `Consumer` method to process the selected messages one by one.
- a `Closer` method to be notified when the connection closes.

The messages selected by Filter are queued until they are processed by the
consumer. Messages are filtered and process in the respective order of their
arrival. Each handler has a queue of size 10.

Improvement: replace the consumer with a channel of message. EndPoint
would send the message to the channel or drop the message if the
buffered channel is full. This to allow of arbitrary buffer size and
more flexible client handling (see improvement of client proxy
handler).

## Client proxy handler

Call data flow: message is sent to the endpoint. reply message is read
by endpoint, sent to the consumer queue, extracted from the queue and
sent to the reply channel (client.go). Then read from the reply
channel, deserialized by the specialized proxy and returned to the
caller.

Subscribe data flow: event message read by endpoint, sent to the
consumer queue, extracted from the queue and sent to the event channel
(client.go). Then read from the event channel, deserialized and sent to
the subscriber channel. Then read from the subscriber's channel.

Improvement: no need of two channels to reach the deserialization
stage. See previous improvement.

## Server message dispatch

Each incoming connection generates a new endpoint. Traffic from the various
endpoints is process in parallel.

Incoming messages from an endpoint are sent to the router (see
Router.Receive). The router dispatches the messages to the various
services. Each service dispatch the messages the mailbox of the
object. Each object has its mailbox which is process sequentially.

At this stage the messages have been sent to handler queue of the
endpoing and then to the object mailbox.

This has two consequences:

    - messages from an endpoint are post in order to the mailbox.
    - objects receive one message at a time.
