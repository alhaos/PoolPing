package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alhaos/PoolPing/internal/config"
	"github.com/alhaos/PoolPing/internal/pingpool"
)

func main() {

	// Init config
	path := flag.String("config", "config/config.yml", "PoolPing config file path")
	flag.Parse()

	conf, err := config.Load(*path)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jobs := make(chan string)

	go func() {
		for _, host := range conf.Hosts {
			jobs <- host
		}
		close(jobs)
	}()

	timeoutDuration := time.Duration(conf.TimeoutMs) * time.Millisecond

	results := pingpool.PingPool(ctx, conf.Workers, timeoutDuration, jobs)

	var aliveCounter, deadCounter int

	fmt.Printf("Results:\n")
	for result := range results {
		if result.Err == nil {
			fmt.Printf("  - %-50s %v\n", result.Host, result.Latency.Round(time.Millisecond))
			aliveCounter++
		} else {
			fmt.Printf("  - %-50s %s\n", result.Host, result.Err)
			deadCounter++
		}
	}
	fmt.Printf("Total host processed: %d\n", len(conf.Hosts))
	fmt.Printf("  - Alive: %d\n", aliveCounter)
	fmt.Printf("  - Dead: %d\n", deadCounter)
}
