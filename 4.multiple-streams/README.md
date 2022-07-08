# 4.Multiple streams
In the last example, they focus on how to establish a quic connection and what packets need to be sent during the connection establishing process. In 
chapter4 here, let's start to find out multiple streams handling.

## Client usage
In this project, client supports to make many stream requests within the same connection. You can press any key except `q`
to continue and enter `q` to quit.

## Things to explore
To be honest, quic is similar to tcp,  as both of them are in transport layer. When we use a tcp based application, 
all we should pay attention to is how to handle the stream based protocol to application protocol.  

In this chapter, we only shortly answer some questions raised when run the demo code. In the following chapters, I think
I should focus on what's the new feature provided by quic, what's the superiorities of quic and what's the suitable cases for quic.
