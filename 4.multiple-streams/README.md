# 4.Multiple streams
In the last examples, they focus on how to establish a quic connection and what's the packets during establishing. In 
chapter4 here, the project starts to do multiple streams handling.

## Client usage
In this project, client side supports to make many stream request with a same connection. Can press any key except `q`
to continue and enter `q` to quit.

## Things to explore
### All streamId could be divided by 4 fully
Investigate why all stream IDs has the relationship `id%4 == 0` and how it is generated.

//TODO

### Cannot get the connection ID
The quic could not get the connection ID, investigate more about the connection ID.  

//TODO