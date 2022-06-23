# quic-server-instance
This chapter focus on how the server instance works.

## Server workflow
The server is designed and encapsulated well. It uses channel to pass data from different go routines which benefits 
from the mechanism of go routine.  
