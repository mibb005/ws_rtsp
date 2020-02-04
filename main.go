package main

import (
	"./rtsp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// go rtsp.StartStream("rtsp://admin:thrn.com@192.168.2.214:554/h264/ch33/main/av_stream")
	go rtsp.StartHttp()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()
	log.Println("Server Start Awaiting Signal")
	<-done
	log.Println("Exiting")
}
