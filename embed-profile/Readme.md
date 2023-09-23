The purpose of this experiment is to understand how the use of `embed.FS`
affects a program's CPU and memory usage as compared to ones that directly
serve the same content from the file system (through streaming or through
a cached instance).

The applications all serve (`GET /tiff`) a base64 string that represents the
[Free_Test_Data_2.3MB_TIFF.tif](./Free_Test_Data_2.3MB_TIFF.tif) acquired
from [https://freetestdata.com/image-files/tiff/](https://freetestdata.com/image-files/tiff/).

## Building

There are three applications:

1. `embedded`
2. `streaming`
3. `cached`

Each should be built by:

```sh
$ go build ./cmd/<name>
```

## Running A Test

The applications are built with internal profiling enabled. We run a test by
starting a command in a terminal like so (e.g. for the `embedded` test):

```sh
$ ./embedded -cpuprofile=embedded_cpu.prof -memprofile=embedded_mem.prof
```

In another terminal, we need to exercise the application by using some HTTP
benchmarking tool, e.g. [autocannon](https://www.npmjs.com/package/autocannon):

```sh
$ autocannon -d 60 -c 30 -w 3 127.0.0.1:8080/tiff
```

After the benchmarking tool has completed, and before we terminate the
application, we can use the included [get_mem.sh](./get_mem.sh) to see how much
memory the process itself used for the test:

```sh
$ ./get_mem.sh 12345 # where `12345` is the pid of the application
```

## Results

See [https://go.dev/blog/pprof](https://go.dev/blog/pprof) for a detailed
discussion on `go tool pprof`.

All results below were generated on a 2021 Apple MacBook Pro with an M1 Max
and 32GB of RAM running macOS 13.4.1. The Go version used was 1.21.1.

### embedded

#### Memory

```sh
$ autocannon -d 60 -c 30 -w 3 127.0.0.1:8080/tiff
Running 60s test @ http://127.0.0.1:8080/tiff
30 connections
3 workers

\
┌─────────┬───────┬───────┬───────┬───────┬──────────┬─────────┬────────┐
│ Stat    │ 2.5%  │ 50%   │ 97.5% │ 99%   │ Avg      │ Stdev   │ Max    │
├─────────┼───────┼───────┼───────┼───────┼──────────┼─────────┼────────┤
│ Latency │ 17 ms │ 26 ms │ 43 ms │ 47 ms │ 27.39 ms │ 6.92 ms │ 121 ms │
└─────────┴───────┴───────┴───────┴───────┴──────────┴─────────┴────────┘
┌───────────┬─────────┬─────────┬─────────┬────────┬─────────┬─────────┬─────────┐
│ Stat      │ 1%      │ 2.5%    │ 50%     │ 97.5%  │ Avg     │ Stdev   │ Min     │
├───────────┼─────────┼─────────┼─────────┼────────┼─────────┼─────────┼─────────┤
│ Req/Sec   │ 1000    │ 1042    │ 1083    │ 1120   │ 1082.85 │ 19.87   │ 1000    │
├───────────┼─────────┼─────────┼─────────┼────────┼─────────┼─────────┼─────────┤
│ Bytes/Sec │ 3.04 GB │ 3.16 GB │ 3.29 GB │ 3.4 GB │ 3.29 GB │ 60.3 MB │ 3.04 GB │
└───────────┴─────────┴─────────┴─────────┴────────┴─────────┴─────────┴─────────┘

Req/Bytes counts sampled once per second.
# of samples: 180

65k requests in 60.44s, 197 GB read
```

```sh
$ ./get_mem.sh <pid>
   RSS      VSZ COMMAND
108M 400347M ./embedded -cpuprofile=embedded_cpu.prof -memprofile=embedded_mem.prof
```

```sh
$ go tool pprof -text embedded_mem.prof 
Type: alloc_space
Time: Sep 23, 2023 at 10:50am (EDT)
Showing nodes accounting for 184.02GB, 99.86% of 184.28GB total
Dropped 49 nodes (cum <= 0.92GB)
      flat  flat%   sum%        cum   cum%
  184.02GB 99.86% 99.86%   184.02GB 99.86%  embed.FS.ReadFile
         0     0% 99.86%   184.08GB 99.89%  main.routeHandler
         0     0% 99.86%   184.08GB 99.89%  net/http.(*ServeMux).ServeHTTP
         0     0% 99.86%   184.27GB   100%  net/http.(*conn).serve
         0     0% 99.86%   184.08GB 99.89%  net/http.HandlerFunc.ServeHTTP
         0     0% 99.86%   184.08GB 99.89%  net/http.serverHandler.ServeHTTP
```

#### CPU

```sh
$ go tool pprof -text embedded_cpu.prof
Type: cpu
Time: Sep 23, 2023 at 10:48am (EDT)
Duration: 84.83s, Total samples = 81.04s (95.53%)
Showing nodes accounting for 77.61s, 95.77% of 81.04s total
Dropped 315 nodes (cum <= 0.41s)
      flat  flat%   sum%        cum   cum%
    42.45s 52.38% 52.38%     42.48s 52.42%  syscall.syscall
     8.26s 10.19% 62.57%      8.26s 10.19%  runtime.madvise
     6.49s  8.01% 70.58%      6.49s  8.01%  runtime.kevent
     4.38s  5.40% 75.99%      4.38s  5.40%  runtime.memmove
     3.59s  4.43% 80.42%      3.59s  4.43%  runtime.pthread_kill
     3.48s  4.29% 84.71%      3.48s  4.29%  runtime.pthread_cond_wait
     2.92s  3.60% 88.31%      2.92s  3.60%  runtime.pthread_cond_signal
     2.09s  2.58% 90.89%      2.09s  2.58%  runtime.usleep
     0.77s  0.95% 91.84%      2.15s  2.65%  runtime.scanobject
     0.69s  0.85% 92.69%      0.69s  0.85%  runtime.pthread_cond_timedwait_relative_np
     0.57s   0.7% 93.40%      0.76s  0.94%  runtime.greyobject
     0.56s  0.69% 94.09%      7.75s  9.56%  runtime.gcDrain
     0.43s  0.53% 94.62%      0.43s  0.53%  runtime.heapBits.nextFast (inline)
     0.28s  0.35% 94.97%      0.58s  0.72%  runtime.scanblock
     0.21s  0.26% 95.22%      0.50s  0.62%  runtime.pcvalue
     0.07s 0.086% 95.31%      0.58s  0.72%  runtime.scanframeworker
     0.05s 0.062% 95.37%      3.01s  3.71%  runtime.markroot
     0.04s 0.049% 95.42%      0.55s  0.68%  runtime.gcDrainN
     0.04s 0.049% 95.47%      0.45s  0.56%  runtime.sweepone
     0.03s 0.037% 95.51%      1.21s  1.49%  runtime.lock2
     0.03s 0.037% 95.55%      0.80s  0.99%  runtime.mallocgc
     0.03s 0.037% 95.58%     12.33s 15.21%  runtime.schedule
     0.02s 0.025% 95.61%        11s 13.57%  runtime.findRunnable
     0.02s 0.025% 95.63%      0.74s  0.91%  runtime.stealWork
     0.01s 0.012% 95.64%     45.54s 56.19%  net/http.(*ServeMux).ServeHTTP
     0.01s 0.012% 95.66%      2.03s  2.50%  runtime.(*gcControllerState).enlistWorker
     0.01s 0.012% 95.67%      2.18s  2.69%  runtime.(*gcWork).balance
     0.01s 0.012% 95.68%      5.24s  6.47%  runtime.(*mheap).allocSpan
     0.01s 0.012% 95.69%      0.44s  0.54%  runtime.bgsweep
     0.01s 0.012% 95.71%      1.14s  1.41%  runtime.forEachP
     0.01s 0.012% 95.72%      3.38s  4.17%  runtime.notesleep
     0.01s 0.012% 95.73%      2.63s  3.25%  runtime.preemptone
     0.01s 0.012% 95.74%      1.34s  1.65%  runtime.suspendG
     0.01s 0.012% 95.76%     20.28s 25.02%  runtime.systemstack
     0.01s 0.012% 95.77%      1.71s  2.11%  runtime.wakep
         0     0% 95.77%      1.12s  1.38%  bufio.(*Reader).Peek
         0     0% 95.77%      1.12s  1.38%  bufio.(*Reader).fill
         0     0% 95.77%         5s  6.17%  bufio.(*Writer).Flush
         0     0% 95.77%     40.59s 50.09%  bufio.(*Writer).Write
         0     0% 95.77%      4.93s  6.08%  embed.FS.ReadFile
         0     0% 95.77%      1.55s  1.91%  internal/poll.(*FD).Read
         0     0% 95.77%     40.93s 50.51%  internal/poll.(*FD).Write
         0     0% 95.77%     42.47s 52.41%  internal/poll.ignoringEINTRIO (inline)
         0     0% 95.77%     45.52s 56.17%  main.routeHandler
         0     0% 95.77%      1.57s  1.94%  net.(*conn).Read
         0     0% 95.77%     40.93s 50.51%  net.(*conn).Write
         0     0% 95.77%      1.55s  1.91%  net.(*netFD).Read
         0     0% 95.77%     40.93s 50.51%  net.(*netFD).Write
         0     0% 95.77%     40.58s 50.07%  net/http.(*chunkWriter).Write
         0     0% 95.77%     47.60s 58.74%  net/http.(*conn).serve
         0     0% 95.77%      1.12s  1.38%  net/http.(*connReader).Read
         0     0% 95.77%      0.46s  0.57%  net/http.(*connReader).backgroundRead
         0     0% 95.77%     40.59s 50.09%  net/http.(*response).Write
         0     0% 95.77%      0.70s  0.86%  net/http.(*response).finishRequest
         0     0% 95.77%     40.59s 50.09%  net/http.(*response).write
         0     0% 95.77%     45.52s 56.17%  net/http.HandlerFunc.ServeHTTP
         0     0% 95.77%     40.93s 50.51%  net/http.checkConnErrorWriter.Write
         0     0% 95.77%     45.54s 56.19%  net/http.serverHandler.ServeHTTP
         0     0% 95.77%      5.08s  6.27%  runtime.(*mheap).alloc.func1
         0     0% 95.77%      3.11s  3.84%  runtime.(*pageAlloc).scavenge.func1
         0     0% 95.77%      3.11s  3.84%  runtime.(*pageAlloc).scavengeOne
         0     0% 95.77%      0.45s  0.56%  runtime.(*unwinder).next
         0     0% 95.77%      0.59s  0.73%  runtime.gcAssistAlloc.func1
         0     0% 95.77%      0.59s  0.73%  runtime.gcAssistAlloc1
         0     0% 95.77%      4.61s  5.69%  runtime.gcBgMarkWorker
         0     0% 95.77%      7.77s  9.59%  runtime.gcBgMarkWorker.func2
         0     0% 95.77%      0.66s  0.81%  runtime.gcMarkDone.func1
         0     0% 95.77%      0.48s  0.59%  runtime.gcMarkTermination.func4
         0     0% 95.77%      0.46s  0.57%  runtime.gcStart.func1
         0     0% 95.77%      0.51s  0.63%  runtime.gcStart.func3
         0     0% 95.77%      0.60s  0.74%  runtime.gcstopm
         0     0% 95.77%      1.20s  1.48%  runtime.gopreempt_m
         0     0% 95.77%      1.20s  1.48%  runtime.goschedImpl
         0     0% 95.77%      1.21s  1.49%  runtime.lock (inline)
         0     0% 95.77%      1.21s  1.49%  runtime.lockWithRank (inline)
         0     0% 95.77%      3.38s  4.17%  runtime.mPark (inline)
         0     0% 95.77%      2.52s  3.11%  runtime.markroot.func1
         0     0% 95.77%     11.05s 13.64%  runtime.mcall
         0     0% 95.77%      1.55s  1.91%  runtime.morestack
         0     0% 95.77%      6.49s  8.01%  runtime.netpoll
         0     0% 95.77%      0.65s   0.8%  runtime.newproc.func1
         0     0% 95.77%      1.56s  1.92%  runtime.newstack
         0     0% 95.77%      0.69s  0.85%  runtime.notetsleep
         0     0% 95.77%      0.69s  0.85%  runtime.notetsleep_internal
         0     0% 95.77%      2.77s  3.42%  runtime.notewakeup
         0     0% 95.77%      1.42s  1.75%  runtime.osyield (inline)
         0     0% 95.77%     10.86s 13.40%  runtime.park_m
         0     0% 95.77%      3.60s  4.44%  runtime.preemptM
         0     0% 95.77%      0.62s  0.77%  runtime.preemptall
         0     0% 95.77%      0.58s  0.72%  runtime.rawbyteslice
         0     0% 95.77%      0.76s  0.94%  runtime.resetspinning
         0     0% 95.77%      0.77s  0.95%  runtime.runSafePointFn
         0     0% 95.77%      0.67s  0.83%  runtime.runqgrab
         0     0% 95.77%      0.67s  0.83%  runtime.runqsteal
         0     0% 95.77%      1.16s  1.43%  runtime.scanstack
         0     0% 95.77%      4.18s  5.16%  runtime.semasleep
         0     0% 95.77%      2.98s  3.68%  runtime.semawakeup
         0     0% 95.77%      3.59s  4.43%  runtime.signalM (inline)
         0     0% 95.77%      0.81s     1%  runtime.startTheWorldWithSema
         0     0% 95.77%      1.70s  2.10%  runtime.startm
         0     0% 95.77%      0.53s  0.65%  runtime.stopTheWorldWithSema
         0     0% 95.77%      0.49s   0.6%  runtime.stoplockedm
         0     0% 95.77%      2.96s  3.65%  runtime.stopm
         0     0% 95.77%      4.93s  6.08%  runtime.stringtoslicebyte
         0     0% 95.77%      3.10s  3.83%  runtime.sysUnused (inline)
         0     0% 95.77%      3.10s  3.83%  runtime.sysUnusedOS (inline)
         0     0% 95.77%      5.16s  6.37%  runtime.sysUsed (inline)
         0     0% 95.77%      5.16s  6.37%  runtime.sysUsedOS (inline)
         0     0% 95.77%      1.55s  1.91%  syscall.Read (inline)
         0     0% 95.77%     40.92s 50.49%  syscall.Write (inline)
         0     0% 95.77%      1.55s  1.91%  syscall.read
         0     0% 95.77%     40.92s 50.49%  syscall.write
```

### streaming

#### Memory

```sh
$ autocannon -d 60 -c 30 -w 3 127.0.0.1:8080/tiff
Running 60s test @ http://127.0.0.1:8080/tiff
30 connections
3 workers

\
┌─────────┬───────┬───────┬────────┬────────┬──────────┬─────────┬────────┐
│ Stat    │ 2.5%  │ 50%   │ 97.5%  │ 99%    │ Avg      │ Stdev   │ Max    │
├─────────┼───────┼───────┼────────┼────────┼──────────┼─────────┼────────┤
│ Latency │ 40 ms │ 94 ms │ 186 ms │ 208 ms │ 99.04 ms │ 37.9 ms │ 320 ms │
└─────────┴───────┴───────┴────────┴────────┴──────────┴─────────┴────────┘
┌───────────┬────────┬────────┬────────┬────────┬────────┬─────────┬────────┐
│ Stat      │ 1%     │ 2.5%   │ 50%    │ 97.5%  │ Avg    │ Stdev   │ Min    │
├───────────┼────────┼────────┼────────┼────────┼────────┼─────────┼────────┤
│ Req/Sec   │ 269    │ 290    │ 307    │ 322    │ 307.12 │ 9.31    │ 269    │
├───────────┼────────┼────────┼────────┼────────┼────────┼─────────┼────────┤
│ Bytes/Sec │ 820 MB │ 884 MB │ 936 MB │ 981 MB │ 936 MB │ 28.4 MB │ 820 MB │
└───────────┴────────┴────────┴────────┴────────┴────────┴─────────┴────────┘

Req/Bytes counts sampled once per second.
# of samples: 177

18k requests in 60.09s, 55.2 GB read
```

```sh

$ ./get_mem.sh <pid>
   RSS      VSZ COMMAND
18M 400440M ./streaming -cpuprofile=streaming_cpu.prof -memprofile=streaming_mem.prof
```

```sh
$ go tool pprof -text streaming_mem.prof
Type: alloc_space
Time: Sep 23, 2023 at 11:14am (EDT)
Showing nodes accounting for 254.53MB, 98.20% of 259.19MB total
Dropped 26 nodes (cum <= 1.30MB)
      flat  flat%   sum%        cum   cum%
  207.50MB 80.06% 80.06%   208.50MB 80.44%  net/http.(*chunkWriter).Write
   19.52MB  7.53% 87.59%   232.02MB 89.52%  main.routeHandler
    7.50MB  2.89% 90.48%     7.50MB  2.89%  net/textproto.readMIMEHeader
    6.50MB  2.51% 92.99%    15.50MB  5.98%  net/http.readRequest
    5.50MB  2.12% 95.11%    24.01MB  9.26%  net/http.(*conn).readRequest
    3.50MB  1.35% 96.46%     3.50MB  1.35%  os.newFile
    1.50MB  0.58% 97.04%     1.50MB  0.58%  bufio.NewWriterSize (inline)
    1.50MB  0.58% 97.62%     1.50MB  0.58%  net/url.parse
    1.50MB  0.58% 98.20%     1.50MB  0.58%  context.withCancel (inline)
         0     0% 98.20%   208.50MB 80.44%  bufio.(*Writer).Flush
         0     0% 98.20%      208MB 80.25%  bufio.(*Writer).Write
         0     0% 98.20%     1.50MB  0.58%  context.WithCancel
         0     0% 98.20%   232.02MB 89.52%  net/http.(*ServeMux).ServeHTTP
         0     0% 98.20%   257.53MB 99.36%  net/http.(*conn).serve
         0     0% 98.20%      208MB 80.25%  net/http.(*response).Write
         0     0% 98.20%      208MB 80.25%  net/http.(*response).write
         0     0% 98.20%   232.02MB 89.52%  net/http.HandlerFunc.ServeHTTP
         0     0% 98.20%     1.50MB  0.58%  net/http.newBufioWriterSize
         0     0% 98.20%   232.02MB 89.52%  net/http.serverHandler.ServeHTTP
         0     0% 98.20%     7.50MB  2.89%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0% 98.20%     1.50MB  0.58%  net/url.ParseRequestURI
         0     0% 98.20%     4.50MB  1.74%  os.Open (inline)
         0     0% 98.20%     4.50MB  1.74%  os.OpenFile
         0     0% 98.20%     4.50MB  1.74%  os.openFileNolog
```

#### CPU

```sh
$ go tool pprof -text streaming_cpu.prof
Type: cpu
Time: Sep 23, 2023 at 11:13am (EDT)
Duration: 84.19s, Total samples = 390.13s (463.41%)
Showing nodes accounting for 388.02s, 99.46% of 390.13s total
Dropped 169 nodes (cum <= 1.95s)
      flat  flat%   sum%        cum   cum%
   387.99s 99.45% 99.45%    388.09s 99.48%  syscall.syscall
     0.02s 0.0051% 99.46%     80.93s 20.74%  bufio.(*Writer).Write
     0.01s 0.0026% 99.46%    305.82s 78.39%  os.(*File).read (inline)
         0     0% 99.46%     81.04s 20.77%  bufio.(*Writer).Flush
         0     0% 99.46%    306.23s 78.49%  internal/poll.(*FD).Read
         0     0% 99.46%     80.85s 20.72%  internal/poll.(*FD).Write
         0     0% 99.46%    387.05s 99.21%  internal/poll.ignoringEINTRIO (inline)
         0     0% 99.46%    387.69s 99.37%  main.routeHandler
         0     0% 99.46%     80.85s 20.72%  net.(*conn).Write
         0     0% 99.46%     80.85s 20.72%  net.(*netFD).Write
         0     0% 99.46%    387.70s 99.38%  net/http.(*ServeMux).ServeHTTP
         0     0% 99.46%     80.90s 20.74%  net/http.(*chunkWriter).Write
         0     0% 99.46%    388.32s 99.54%  net/http.(*conn).serve
         0     0% 99.46%     80.83s 20.72%  net/http.(*response).Write
         0     0% 99.46%     80.83s 20.72%  net/http.(*response).write
         0     0% 99.46%    387.69s 99.37%  net/http.HandlerFunc.ServeHTTP
         0     0% 99.46%     80.85s 20.72%  net/http.checkConnErrorWriter.Write
         0     0% 99.46%    387.70s 99.38%  net/http.serverHandler.ServeHTTP
         0     0% 99.46%    305.82s 78.39%  os.(*File).Read
         0     0% 99.46%    306.20s 78.49%  syscall.Read (inline)
         0     0% 99.46%     80.85s 20.72%  syscall.Write (inline)
         0     0% 99.46%    306.20s 78.49%  syscall.read
         0     0% 99.46%     80.85s 20.72%  syscall.write
```

### cached

#### Memory

```sh
$ autocannon -d 60 -c 30 -w 3 127.0.0.1:8080/tiff
Running 60s test @ http://127.0.0.1:8080/tiff
30 connections
3 workers

/
┌─────────┬───────┬───────┬───────┬───────┬──────────┬─────────┬───────┐
│ Stat    │ 2.5%  │ 50%   │ 97.5% │ 99%   │ Avg      │ Stdev   │ Max   │
├─────────┼───────┼───────┼───────┼───────┼──────────┼─────────┼───────┤
│ Latency │ 17 ms │ 26 ms │ 43 ms │ 47 ms │ 27.07 ms │ 6.81 ms │ 83 ms │
└─────────┴───────┴───────┴───────┴───────┴──────────┴─────────┴───────┘
┌───────────┬─────────┬─────────┬─────────┬─────────┬─────────┬────────┬─────────┐
│ Stat      │ 1%      │ 2.5%    │ 50%     │ 97.5%   │ Avg     │ Stdev  │ Min     │
├───────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────┼─────────┤
│ Req/Sec   │ 951     │ 997     │ 1105    │ 1160    │ 1095.2  │ 48.85  │ 951     │
├───────────┼─────────┼─────────┼─────────┼─────────┼─────────┼────────┼─────────┤
│ Bytes/Sec │ 2.89 GB │ 3.03 GB │ 3.36 GB │ 3.52 GB │ 3.33 GB │ 148 MB │ 2.89 GB │
└───────────┴─────────┴─────────┴─────────┴─────────┴─────────┴────────┴─────────┘

Req/Bytes counts sampled once per second.
# of samples: 180

66k requests in 60.38s, 200 GB read
```

```sh 
$ ./get_mem.sh <pid>
   RSS      VSZ COMMAND
137M 400376M ./cached -cpuprofile=cached_cpu.prof -memprofile=cached_mem.prof
```

```sh
$ go tool pprof -text cached_mem.prof
Type: alloc_space
Time: Sep 23, 2023 at 11:28am (EDT)
Showing nodes accounting for 406.45MB, 98.51% of 412.61MB total
Dropped 12 nodes (cum <= 2.06MB)
      flat  flat%   sum%        cum   cum%
  314.93MB 76.33% 76.33%   314.93MB 76.33%  io.ReadAll
   31.01MB  7.52% 83.84%    31.01MB  7.52%  net/textproto.readMIMEHeader
   26.50MB  6.42% 90.27%    89.02MB 21.57%  net/http.(*conn).readRequest
      11MB  2.67% 92.93%    55.01MB 13.33%  net/http.readRequest
      11MB  2.67% 95.60%       11MB  2.67%  net/url.parse
    5.50MB  1.33% 96.93%     5.50MB  1.33%  context.withCancel (inline)
       5MB  1.21% 98.14%        5MB  1.21%  net.(*conn).Read
    1.50MB  0.36% 98.51%        7MB  1.70%  context.WithCancel
         0     0% 98.51%   315.93MB 76.57%  main.routeHandler
         0     0% 98.51%   315.93MB 76.57%  net/http.(*ServeMux).ServeHTTP
         0     0% 98.51%   406.45MB 98.51%  net/http.(*conn).serve
         0     0% 98.51%        5MB  1.21%  net/http.(*connReader).backgroundRead
         0     0% 98.51%   315.93MB 76.57%  net/http.HandlerFunc.ServeHTTP
         0     0% 98.51%   315.93MB 76.57%  net/http.serverHandler.ServeHTTP
         0     0% 98.51%    31.01MB  7.52%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0% 98.51%       11MB  2.67%  net/url.ParseRequestURI
```

#### CPU

```sh
$ go tool pprof -text cached_cpu.prof
Type: cpu
Time: Sep 23, 2023 at 11:26am (EDT)
Duration: 74.21s, Total samples = 42.03s (56.64%)
Showing nodes accounting for 41.54s, 98.83% of 42.03s total
Dropped 109 nodes (cum <= 0.21s)
      flat  flat%   sum%        cum   cum%
    35.99s 85.63% 85.63%     35.99s 85.63%  syscall.syscall
     3.17s  7.54% 93.17%      3.17s  7.54%  runtime.kevent
     1.25s  2.97% 96.15%      1.25s  2.97%  runtime.pthread_cond_wait
     1.11s  2.64% 98.79%      1.11s  2.64%  runtime.pthread_cond_signal
     0.01s 0.024% 98.81%      1.47s  3.50%  bufio.(*Writer).Flush
     0.01s 0.024% 98.83%      4.60s 10.94%  runtime.findRunnable
         0     0% 98.83%      0.49s  1.17%  bufio.(*Reader).Peek
         0     0% 98.83%      0.49s  1.17%  bufio.(*Reader).fill
         0     0% 98.83%     34.82s 82.85%  bufio.(*Writer).Write
         0     0% 98.83%      0.71s  1.69%  internal/poll.(*FD).Read
         0     0% 98.83%     35.29s 83.96%  internal/poll.(*FD).Write
         0     0% 98.83%     35.99s 85.63%  internal/poll.ignoringEINTRIO (inline)
         0     0% 98.83%     34.88s 82.99%  main.routeHandler
         0     0% 98.83%      0.68s  1.62%  net.(*conn).Read
         0     0% 98.83%     35.29s 83.96%  net.(*conn).Write
         0     0% 98.83%      0.68s  1.62%  net.(*netFD).Read
         0     0% 98.83%     35.29s 83.96%  net.(*netFD).Write
         0     0% 98.83%     34.89s 83.01%  net/http.(*ServeMux).ServeHTTP
         0     0% 98.83%     34.82s 82.85%  net/http.(*chunkWriter).Write
         0     0% 98.83%     36.05s 85.77%  net/http.(*conn).serve
         0     0% 98.83%      0.49s  1.17%  net/http.(*connReader).Read
         0     0% 98.83%     34.82s 82.85%  net/http.(*response).Write
         0     0% 98.83%      0.51s  1.21%  net/http.(*response).finishRequest
         0     0% 98.83%     34.82s 82.85%  net/http.(*response).write
         0     0% 98.83%     34.88s 82.99%  net/http.HandlerFunc.ServeHTTP
         0     0% 98.83%     35.29s 83.96%  net/http.checkConnErrorWriter.Write
         0     0% 98.83%     34.89s 83.01%  net/http.serverHandler.ServeHTTP
         0     0% 98.83%      1.25s  2.97%  runtime.mPark (inline)
         0     0% 98.83%      4.87s 11.59%  runtime.mcall
         0     0% 98.83%      3.19s  7.59%  runtime.netpoll
         0     0% 98.83%      0.69s  1.64%  runtime.newproc.func1
         0     0% 98.83%      1.25s  2.97%  runtime.notesleep
         0     0% 98.83%      1.15s  2.74%  runtime.notewakeup
         0     0% 98.83%      4.86s 11.56%  runtime.park_m
         0     0% 98.83%      0.25s  0.59%  runtime.resetspinning
         0     0% 98.83%      4.87s 11.59%  runtime.schedule
         0     0% 98.83%      1.25s  2.97%  runtime.semasleep
         0     0% 98.83%      1.15s  2.74%  runtime.semawakeup
         0     0% 98.83%      1.05s  2.50%  runtime.startm
         0     0% 98.83%      1.26s  3.00%  runtime.stopm
         0     0% 98.83%      0.94s  2.24%  runtime.systemstack
         0     0% 98.83%      1.05s  2.50%  runtime.wakep
         0     0% 98.83%      0.70s  1.67%  syscall.Read (inline)
         0     0% 98.83%     35.29s 83.96%  syscall.Write (inline)
         0     0% 98.83%      0.70s  1.67%  syscall.read
         0     0% 98.83%     35.29s 83.96%  syscall.write
```
