package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/posener/h2conn"
	"golang.org/x/net/http2"
)

const (
	serverURL1 = "https://http2.golang.org/ECHO"
	serverURL2 = "https://localhost:8000"
)

// A client for Go's HTTP2 echo server example at http2.golang.org/ECHO

func main() {
	// Create a client, that uses the HTTP PUT method.
	c := h2conn.Client{
		Method: http.MethodGet,
		Client: &http.Client{Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}},
	}

	// Connect to the HTTP2 server
	// The returned conn can be used to:
	//   1. Write - send data to the server.
	//   2. Read - receive data from the server.
	conn, resp, err := c.Connect(context.Background(), serverURL1)
	if err != nil {
		log.Fatal(err, "xxx")
	}
	defer conn.Close()
	log.Printf("Got: %d", resp.StatusCode)

	// Send time periodically to the server
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Fprintf(conn, "It is now %v\n", time.Now())
		}
	}()

	// Read responses from the server to the stdout.
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		log.Fatal(err)
	}
}
