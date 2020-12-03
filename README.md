# Protocols
* [Transport layer](#Transport layer)
* [Security layer](#Security layer)

In this project I've implemented two layers of a protocol, the transport layer atop udp and the security layer.


## Transport layer:

I've build a transport layer atod UDP protocol;
I've used some possibilities which golang have, something like functions which listens on a port, dial to a port. Also the write function and the readFromUDP function which helped me to share information between my client and server.This native golang functions was wraped in my own own functions;

  My transport layer has function validation, which checks if the checksums of the messages which go from client to server are valid or not and if they are valid my server creates a new packet with "ack"(or "nack") and send it back to the client to inform if the data arrive and the connection was established succesfully.
  Also my client could send message written by hand, which after that could be handled by the server;
  
  
  
  
## Security layer: 

  In the security layer I've implemented the Diffie hellman key exchange. Here I've use two public keys, and respectivelly two private keys.After the performing some operation my client receive a message which say that the secured connection was established successfully. After this operation the client and the server receive two shared key which in future is used to encrypt and decrypt data;
  
  The calculation are performed in the next way:
      
  ![photo5348418513833930471](https://user-images.githubusercontent.com/47230162/101048175-18ef6300-358b-11eb-80d6-50a67b3e46e4.jpg)

  After all of this magic stuff happens our client and server have their shared keys and can encryot and decrypt data using symetric encryption;
  
 
 
 Here is the link to drive where I upload the video: https://drive.google.com/file/d/104noIk0BW8v6TdP_oEbQtECpXc_DYKEh/view?usp=sharing
  
  
  
  
