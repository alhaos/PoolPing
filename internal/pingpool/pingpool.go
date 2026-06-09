// Package pingpool implements a worker pool for parallel ICMP ping checks.
// It uses github.com/prometheus-community/pro-bing under the hood,
// which does not require elevated privileges on Windows.
package pingpool

import (
	"context"
	"fmt"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

// Result holds the outcome of a single ping check.
type Result struct {
	// Host is the target hostname or IP address that was pinged.
	Host string

	// Latency is the round-trip time of the ICMP echo request.
	// Zero if the host did not reply.
	Latency time.Duration

	// Err is non-nil if the host was unreachable or an error occurred.
	Err error
}

// PingPool starts numWorkers goroutines that read hostnames from jobs,
// ping each one with the given timeout, and send results to the returned channel.
// The returned channel is closed when all jobs are processed or ctx is cancelled.
func PingPool(ctx context.Context, numWorkers int, timeout time.Duration, jobs <-chan string) <-chan Result {
	// Init output channel
	out := make(chan Result)
	// Init WaitGroup
	var wg sync.WaitGroup
	// Run workers
	for range numWorkers {
		wg.Go(func() {
			// Get job from channel
			for job := range jobs {
				select {
				// Wait until someone reads from out
				case out <- ping(ctx, job, timeout):
				// or context is done
				case <-ctx.Done():
					return
				}
			}
		})
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ping sends a single ICMP echo request to host and returns the result.
// It resolves the hostname, sends one packet, and waits up to timeout for a reply.
func ping(ctx context.Context, host string, timeout time.Duration) Result {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return Result{Host: host, Err: err}
	}

	pinger.Count = 1
	pinger.Timeout = timeout
	pinger.SetPrivileged(true) // required on Windows

	if err := pinger.RunWithContext(ctx); err != nil {
		return Result{Host: host, Err: err}
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return Result{Host: host, Err: fmt.Errorf("no reply")}
	}

	return Result{Host: host, Latency: stats.AvgRtt}
}
