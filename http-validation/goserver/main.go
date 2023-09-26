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

	validate "github.com/go-playground/validator/v10"
	"github.com/jsumners/go-rfc3339"
)

var validator *validate.Validate

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
	validator = validate.New(validate.WithRequiredStructEnabled())
	return validator.RegisterValidation("dateTime", func(fl validate.FieldLevel) bool {
		return rfc3339.IsDateTimeString(fl.Field().String())
	})
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
	Subject string `json:"subject" validate:"required"`
	Created string `json:"created" validate:"required,dateTime"`
	Body    string `json:"body" validate:"required"`
	Author  struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"email"`
	} `json:"author" validate:"required"`
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

	err = validator.Struct(body)
	if err != nil {
		http.Error(res, fmt.Sprintf("validation failed: %v", err), http.StatusBadRequest)
		return
	}

	io.WriteString(res, "ok")
}
