package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	stderrLogger "log"
	stdoutLogger "log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	udpListenAddress  = flag.String("udp-listen-address", "0.0.0.0:9090", "The address to listen on for udp stats requests.")
	httpListenAddress = flag.String("http-listen-address", "0.0.0.0:9091", "The address to listen on for http stat and telemetry requests.")
)

func main() {
	flag.Parse()

	stderrLogger.SetOutput(os.Stderr)
	stdoutLogger.SetOutput(os.Stdout)

	errrorLogger := stderrLogger.Default()
	infoLogger := stdoutLogger.Default()

	infoLogger.Print("welcome to messagebird/pushprom")

	ctx, cancel := context.WithCancel(context.Background())

	var err error

	*udpListenAddress, err = ListenAddress(*udpListenAddress)
	if err != nil {
		errrorLogger.Fatalf(err.Error())
	}

	*httpListenAddress, err = ListenAddress(*httpListenAddress)
	if err != nil {
		errrorLogger.Fatalf(err.Error())
	}

	infoLogger.Print("starting listeners")

	var wg sync.WaitGroup

	go listenUDP(&wg, ctx, errrorLogger, infoLogger)
	go listenHTTP(&wg, ctx, errrorLogger, infoLogger)
	wg.Add(2)

	handleSIGTERM(&wg, cancel, infoLogger)
}

func handleSIGTERM(wg *sync.WaitGroup, cancel func(), infoLogger *log.Logger) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt)
	<-sigc
	infoLogger.Println("received termination/interrupt signal, will terminate")
	cancel()
	defer close(sigc)

	infoLogger.Println("waiting for all servers to gracefully terminate")
	wg.Wait()
	time.Sleep(1 * time.Second)
	infoLogger.Println("all goroutines gracefully finished")
}

// ListenAddress Format a correct listen address
func ListenAddress(s string) (string, error) {
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}
