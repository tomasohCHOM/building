# HTTP from TCP

Code I wrote following ThePrimeagen's [HTTP from TCP](https://www.boot.dev/courses/learn-http-protocol-golang) course on [Boot.dev](https://www.boot.dev/)

## Key Learnings

I'm glad I went through the entire course since I learned a bunch of new stuff from it. Things I genuinely thought were "magic" now seem to be simpler than they really are, and writing HTTP from scratch once gave me a lot of confidence to tackle other challenges as well.

Some of my key takeaways & learnings:

- Pull vs push model: when it comes to listening to data from network connections, it's very similar to the idea of reading from a file. The key difference is you often "pull" data from a file (that is, you read data whenever you want), whereas data is "pushed" to you from a network connection. The server must be ready to handle the data that gets pushed from some sort of client so that it can respond accordingly (this might be somewhat wrong, but it's the way I interpret it)
- HTTP semantics are defined under the RFC 9110, with RFC 9112 including information about HTTP/1.1 (the protocol that was built in the course)

- HTTP requests & responses operate on HTTP messages, following a standard format that is composed by:
    - a start line (e.g. GET / HTTP/1.1)
    - headers (metadata about the message e.g. `Content-Type: application/json`)
    - empty line indicating end of headers (with CRLF)
    - optional body containing data payload
- There are several versions of HTTP (1.1, 2, 3), each having basically the same semantics but implemented in different ways:
    - HTTP 2 does things that HTTP 1.1 like stream multiplexing and header compression for a more efficient protocol
    - HTTP 3 is built on top of QUIC, which runs on UDP instead of TCP
- CRLFs, very important
- Building a state machine to parse an HTTP request, going through each part (start line, headers, body, etc.) one by one while using a buffer to read data to

Genuinely a fun time and it's amazing that the content itself is completely free.

*NOTE*: this isn't meant to be a fully correct implementation of HTTP at all, it was mainly for the learning experience and having a basic understanding of how things work under the hood.

## Resources

- [Boot.dev](https://www.boot.dev/)
- [HTTP from TCP Course](https://www.boot.dev/courses/learn-http-protocol-golang)
- [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html) && [RFC 9112](https://www.rfc-editor.org/rfc/rfc9112.html)
- [Solutions repo I used](https://github.com/oleshko-g/httpfromtcp)
