package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/edlundin/enocean-esp3/pkg"
)

func main() {
	ports, err := pkg.GetSerialPortList()

	if err != nil {
		fmt.Println(fmt.Errorf("unable to get serial port list: %w", err))
	}

	if len(ports) == 0 {
		fmt.Println("no serial ports found")
	}

	ctx, termination := context.WithCancel(context.Background())

	p, channels, err := pkg.OpenSerialPort(ctx, "/dev/cu.usbserial-EO7BBSFR")

	if err != nil {
		fmt.Println(fmt.Errorf("unable to open serial port: %w", err))
		return
	}

	go func() {
		for msg := range channels.All {
			fmt.Println(msg.Kind, msg.Data, msg.Err)
		}
	}()

	terminationChannel := make(chan os.Signal)
	signal.Notify(terminationChannel, syscall.SIGINT, syscall.SIGTERM)
	<-terminationChannel // Blocks here until either SIGINT or SIGTERM is received.

	termination()

	defer p.Close()
}
