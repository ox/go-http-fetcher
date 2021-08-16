# Go HTTP Fetcher

This program watches for `.req` files created in a directory specified by the `-dir` flag, performs an HTTP request, and writes the response out to a `.res` file.

#### Example

```
$ mkdir op
$ go-http-fetcher -dir op &
$ echo "GET http://www.example.com" >> op/example.req
200 GET http://www.example.com
```

## Installation

```
go install github.com/ox/go-http-fetcher
```

## Request File Format

The request file is trimmed of whitespace and then split by spaces into at most 3 items. The first item is the request method, the second is the URL, the optional third item is the request body.

Some valid request bodies:

```
GET http://www.example.com
```

```
POST https://reqbin.com/echo/post/json {"hello": "world"}
```
