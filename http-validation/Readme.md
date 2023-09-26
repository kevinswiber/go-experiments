The purpose of this experiment is to determine how an HTTP body validation
strategy compares to Fastify's standard input validation.

## Building

The tests are run in a Docker container that limits the CPU availability to
1 CPU. As a convenience, a [Taskfile](https://taskfile.dev) is included to
facilitate building things:

```sh
$ task dockerize
```

## Running Tests

The test script, and one of the test servers, requires `npm install` be run.
Subsequent to that, we can run each test:

### Go Server

1. In one terminal: `task run-go`
2. In another terminal: `node test.js`
3. `docker kill go_test_go`

### Node Server

1. In one terminal: `task run-node`
2. In another terminal: `node test.js`
3. `docker kill go_test_node`


## Results

All results below were generated on a 2021 Apple MacBook Pro with an M1 Max
and 32GB of RAM running macOS 13.4.1. The Go version used was 1.21.1. The
Node version used was 20 (via the Docker image `node:20-slim`).

### Go Server

```sh
running valid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬────────┬─────┬─────┬─────┐
│ (index)  │ 1xx │  2xx   │ 3xx │ 4xx │ 5xx │
├──────────┼─────┼────────┼─────┼─────┼─────┤
│ statuses │  0  │ 675266 │  0  │  0  │  0  │
└──────────┴─────┴────────┴─────┴─────┴─────┘
┌──────────┬──────────┬──────────┬─────────┬────────┐
│ (index)  │ average  │   mean   │ stddev  │ total  │
├──────────┼──────────┼──────────┼─────────┼────────┤
│ requests │ 22509.87 │ 22509.87 │ 1174.93 │ 675266 │
└──────────┴──────────┴──────────┴─────────┴────────┘


running invalid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────┬─────┬────────┬─────┐
│ (index)  │ 1xx │ 2xx │ 3xx │  4xx   │ 5xx │
├──────────┼─────┼─────┼─────┼────────┼─────┤
│ statuses │  0  │  0  │  0  │ 596474 │  0  │
└──────────┴─────┴─────┴─────┴────────┴─────┘
┌──────────┬──────────┬──────────┬────────┬────────┐
│ (index)  │ average  │   mean   │ stddev │ total  │
├──────────┼──────────┼──────────┼────────┼────────┤
│ requests │ 19882.14 │ 19882.14 │ 770.46 │ 596474 │
└──────────┴──────────┴──────────┴────────┴────────┘
```

### Node Server

```sh
running valid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬────────┬─────┬─────┬─────┐
│ (index)  │ 1xx │  2xx   │ 3xx │ 4xx │ 5xx │
├──────────┼─────┼────────┼─────┼─────┼─────┤
│ statuses │  0  │ 470603 │  0  │  0  │  0  │
└──────────┴─────┴────────┴─────┴─────┴─────┘
┌──────────┬─────────┬─────────┬─────────┬────────┐
│ (index)  │ average │  mean   │ stddev  │ total  │
├──────────┼─────────┼─────────┼─────────┼────────┤
│ requests │ 15686.8 │ 15686.8 │ 1117.97 │ 470603 │
└──────────┴─────────┴─────────┴─────────┴────────┘


running invalid payload test...
Running 30s test @ http://127.0.0.1:8080/message/new
30 connections
3 workers

┌──────────┬─────┬─────┬─────┬────────┬─────┐
│ (index)  │ 1xx │ 2xx │ 3xx │  4xx   │ 5xx │
├──────────┼─────┼─────┼─────┼────────┼─────┤
│ statuses │  0  │  0  │  0  │ 379683 │  0  │
└──────────┴─────┴─────┴─────┴────────┴─────┘
┌──────────┬──────────┬──────────┬────────┬────────┐
│ (index)  │ average  │   mean   │ stddev │ total  │
├──────────┼──────────┼──────────┼────────┼────────┤
│ requests │ 12657.07 │ 12657.07 │ 517.68 │ 379683 │
└──────────┴──────────┴──────────┴────────┴────────┘
```
