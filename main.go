package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"net"
)

func handle(src net.Conn, host string) {
	dst, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatalln("Unable to connect to out destination host")
	}
	defer dst.Close()

	// Run in goroutine to prevent io.Copy from blocking
	go func() {
		// Copy our source's output to the destination
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()
	// Copy out destination's output back to our source
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

var (
	listenPort *int
	targetURL  string
)

func init() {
	listenPort = flag.Int("p", 8080, "The port on which the proxy will listen")
	flag.StringVar(&targetURL, "u", "localhost:80", "The URL where the proxy should redirect the traffic. Please, use the structure domain:port")
	flag.Parse()
}

func main() {
	// Listen on local port
	portSocket := fmt.Sprintf(":%d", *listenPort)
	lister, err := net.Listen("tcp", portSocket)
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	for {
		conn, err := lister.Accept()
		if err != nil {
			log.Fatal("Unable to accept connection")
		}
		fmt.Println(conn.RemoteAddr())
		go handle(conn, targetURL)
	}

}
