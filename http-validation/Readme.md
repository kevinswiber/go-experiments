The purpose of this experiment is to determine how an HTTP body validation
strategy compares to Fastify's standard input validation.

## Running Tests

The test script, and one of the test servers, requires `npm install` be run.
Subsquent to that, we can run each test:

### Go Server

1. In one terminal: `go run ./goserver/main.go`
2. In another terminal: `node test.js`

### Node Server

1. In one terminal: `node nodeserver/index.js`
2. In another terminal: `node test.js`


## Results

All results below were generated on a 2021 Apple MacBook Pro with an M1 Max
and 32GB of RAM running macOS 13.4.1. The Go version used was 1.21.1.

### Go Server

```sh
running valid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────────┬─────┬─────┬─────┐
│ (index)  │ 1xx │   2xx   │ 3xx │ 4xx │ 5xx │
├──────────┼─────┼─────────┼─────┼─────┼─────┤
│ statuses │  0  │ 3421574 │  0  │  0  │  0  │
└──────────┴─────┴─────────┴─────┴─────┴─────┘
┌──────────┬───────────┬───────────┬─────────┬─────────┐
│ (index)  │  average  │   mean    │ stddev  │  total  │
├──────────┼───────────┼───────────┼─────────┼─────────┤
│ requests │ 114050.14 │ 114050.14 │ 5591.07 │ 3421574 │
└──────────┴───────────┴───────────┴─────────┴─────────┘


running invalid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────┬─────┬─────────┬─────┐
│ (index)  │ 1xx │ 2xx │ 3xx │   4xx   │ 5xx │
├──────────┼─────┼─────┼─────┼─────────┼─────┤
│ statuses │  0  │  0  │  0  │ 3017739 │  0  │
└──────────┴─────┴─────┴─────┴─────────┴─────┘
┌──────────┬──────────┬──────────┬─────────┬─────────┐
│ (index)  │ average  │   mean   │ stddev  │  total  │
├──────────┼──────────┼──────────┼─────────┼─────────┤
│ requests │ 100588.8 │ 100588.8 │ 3323.79 │ 3017739 │
└──────────┴──────────┴──────────┴─────────┴─────────┘
```

### Node Server

```sh
running valid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────────┬─────┬─────┬─────┐
│ (index)  │ 1xx │   2xx   │ 3xx │ 4xx │ 5xx │
├──────────┼─────┼─────────┼─────┼─────┼─────┤
│ statuses │  0  │ 1549706 │  0  │  0  │  0  │
└──────────┴─────┴─────────┴─────┴─────┴─────┘
┌──────────┬──────────┬──────────┬────────┬─────────┐
│ (index)  │ average  │   mean   │ stddev │  total  │
├──────────┼──────────┼──────────┼────────┼─────────┤
│ requests │ 51656.54 │ 51656.54 │ 694.67 │ 1549706 │
└──────────┴──────────┴──────────┴────────┴─────────┘


running invalid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────┬─────┬─────────┬─────┐
│ (index)  │ 1xx │ 2xx │ 3xx │   4xx   │ 5xx │
├──────────┼─────┼─────┼─────┼─────────┼─────┤
│ statuses │  0  │  0  │  0  │ 1102786 │  0  │
└──────────┴─────┴─────┴─────┴─────────┴─────┘
┌──────────┬─────────┬─────────┬────────┬─────────┐
│ (index)  │ average │  mean   │ stddev │  total  │
├──────────┼─────────┼─────────┼────────┼─────────┤
│ requests │ 36758.4 │ 36758.4 │ 265.35 │ 1102786 │
└──────────┴─────────┴─────────┴────────┴─────────┘
```
