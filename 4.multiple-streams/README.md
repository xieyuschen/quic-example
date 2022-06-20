# 4.Multiple streams
In the last examples, they focus on how to establish a quic connection and what's the packets during establishing. In 
chapter4 here, the project starts to do multiple streams handling.

## Client usage
In this project, client side supports to make many stream request with a same connection. Can press any key except `q`
to continue and enter `q` to quit.

## Things to explore
To be honest, quic is similar to tcp as all of them are both in transport layer. When we use a tcp based application, 
all we should pay attention is how to handle the stream based protocol to application protocol.  

In this chapter, we only shortly answer some questions raised when run the demo code. In the following chapters, I think
I should focus on what's the new feature provided by quic, what's the advances of quic and what's the suitable cases for 
the quic.
