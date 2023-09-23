package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	println("pid:", os.Getpid())

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	cpuprofile := flags.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flags.String("memprofile", "", "write mem profile to file")

	err := flags.Parse(args[1:])
	if err != nil {
		return err
	}

	if *cpuprofile != "" {
		cpuFile, err := os.Create(*cpuprofile)
		if err != nil {
			return err
		}
		defer cpuFile.Close()
		err = pprof.StartCPUProfile(cpuFile)
		if err != nil {
			return err
		}
		defer pprof.StopCPUProfile()
	}

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	mux := http.NewServeMux()
	mux.HandleFunc("/tiff", routeHandler)
	server := http.Server{Handler: mux, Addr: "127.0.0.1:8080"}
	go func() {
		server.ListenAndServe()
	}()

mainLoop:
	for {
		select {
		case <-sigCh:
			server.Close()
			break mainLoop
		}
	}

	println("nonCacheHitCount:", nonCacheHitCount)
	println("cacheHitCount:", cacheHitCount)

	if *memprofile != "" {
		memFile, err := os.Create(*memprofile)
		if err != nil {
			return err
		}
		defer memFile.Close()
		runtime.GC()
		err = pprof.Lookup("allocs").WriteTo(memFile, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

var cachedFile *[]byte
var cacheHitCount = 0
var nonCacheHitCount = 0

func routeHandler(res http.ResponseWriter, req *http.Request) {
	if cachedFile != nil {
		cacheHitCount += 1
		res.Write(*cachedFile)
		return
	}

	payload, err := os.Open("resources/test.data")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		io.WriteString(res, err.Error())
		return
	}
	defer payload.Close()

	fileBytes, err := io.ReadAll(bufio.NewReader(payload))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	nonCacheHitCount += 1
	cachedFile = &fileBytes
	res.Write(fileBytes)
}
