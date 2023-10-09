package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xeipuuv/gojsonschema"
)

var validator *gojsonschema.Schema

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	err := initValidator()
	if err != nil {
		return err
	}

	server := initServer()

mainLoop:
	for {
		select {
		case <-sigChan:
			server.Close()
			break mainLoop
		}
	}

	return nil
}

func initValidator() error {
	loader := gojsonschema.NewStringLoader(`
    {
      "$schema": "http://json-schema.org/draft-07/schema#",
      "type": "object",
      "properties": {
        "subject": { "type": "string" },
        "created": { "type": "string", "format": "date-time" },
        "body": { "type": "string" },
        "author": {
          "type": "object",
          "properties": {
            "name": { "type": "string" },
            "email": { "type": "string", "format": "email" }
          },
          "required": ["name", "email"]
        }
      },
      "required": ["subject", "created", "body", "author"]
    }
  `)

	var err error
	validator, err = gojsonschema.NewSchema(loader)
	return err
}

func initServer() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/message/new", routeHandler)

	listenAddr := os.Getenv("HOST_ADDR")
	if listenAddr == "" {
		listenAddr = "127.0.0.1"
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:8080", listenAddr),
		Handler: mux,
	}

	go func() {
		server.ListenAndServe()
	}()

	return &server
}

type postBody struct {
	Subject string `json:"subject"`
	Created string `json:"created"`
	Body    string `json:"body"`
	Author  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

func routeHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	var body postBody
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		http.Error(res, "invalid json body", http.StatusBadRequest)
		return
	}

	loader := gojsonschema.NewGoLoader(body)
	result, err := validator.Validate(loader)
	if err != nil {
		http.Error(res, fmt.Sprintf("validation failed: %v", err), http.StatusBadRequest)
		return
	}
	if !result.Valid() {
		var errors []string = make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			errors = append(errors, fmt.Sprintf("%s", err))
		}
		http.Error(res, fmt.Sprintf("validation failed: %v", errors), http.StatusBadRequest)
		return
	}

	io.WriteString(res, "ok")
}
