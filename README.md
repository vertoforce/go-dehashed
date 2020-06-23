# Dehashed

Go library to use the [dehashed api](https://www.dehashed.com/docs).

## Rate Limits

At the time of writing _"The DeHashed API will only accept 5 Requests / 250ms from a single I.P & API Credential"_.

This is built in to the library and every function call will wait for the next request to be available.
