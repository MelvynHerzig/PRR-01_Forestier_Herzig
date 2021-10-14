# Exercise Protocol Design

## What transport protocol do we use?
We use TCP protocol to communicate between clients and the server

## How does the client find the server (addresses and ports)?
The server will use the localhost adress and the port 8000

## Who speaks first?
The server will speak first.    
A welcome message followed by a list of command will be sent to the client. 


## What is the sequence of messages exchanged by the client and the server? (flow)
S : Welcome in the FH Hotel\
Please identify yourself. Use command "LOGIN {your name}" CRLF\
C : LOGIN John CRLF\
S : Hello John\
You can book our rooms (1-3) for days (1-5)\
Here is the list of commands you can use :
- BOOK roomNumber, arrivalDay, nbNights\
- ROOMLIST day\
- FREEROOM arrivalDay nbNights\
- QUIT CRLF\

C : ROOMLIST 1 CRLF\
S : RESULT 1,"" 2,"" 3"Client1" CRLF\
C : FREEROOM 1 2 CRLF\
S : RESULT 1 2 CRLF\
C : BOOK 1 1 2 CRLF\
S : RESULT OK CRLF\
C : QUIT


## What happens when a message is received from the other party? (semantics)
The server must compute the message from descripted syntax and return :
- The result if possible
- An error otherwise

The client must compute the message from descripty syntax and show :
- The result
- The reason of the error

## What is the syntax of the messages? How we generate and parse them? (syntax)
| Utility | Syntax |
|---|----|
| S'identifier à l'hôtel | LOGIN {name} |
| Book a room (from client) | BOOK {room number} {arrival day} {number of nights}   |
| Get occupations list (from client) | ROOMLIST {day}    |
| Get which room is free for {number of nights} <br>from {arrival day} (from client) | FREEROOM {arrival day} {number of nights}   |
| Close the connection (from client)  | QUIT  |
| Return the result (from server)  | RESULT {resultat}   |
| Return an error (from server) | ERROR {error number} {Reason of the error}   |



## Who closes the connection and when?
The client close the connection using QUIT command
