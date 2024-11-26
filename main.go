package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() { 
	var wg sync.WaitGroup

	wg.Add(1);
	go usuaulRateLimiter(&wg);


	wg.Wait();

	wg.Add(1);
	go burstRateLimiter(&wg);
	wg.Wait();

}

func burstRateLimiter(wg *sync.WaitGroup) { 
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel();
	defer wg.Done();

	burstRateLimiter := make(chan time.Time, 3);

	for i := 0; i  < 3; i++ { 
		burstRateLimiter <- time.Now();
	}

	go func (ctx context.Context) { 
		ticker := time.NewTicker(200 * time.Millisecond);
		defer ticker.Stop()

		for { 
			select { 
				case <- ctx.Done():
					close(burstRateLimiter)
					return
				case t := <-ticker.C: 
					select { 
						case burstRateLimiter <- t:
						default:
					}
			}
		}
	}(ctx)

	burstRequests := make(chan int, 5);

	for i:= range(5) { 
		burstRequests <- i;
	}
	close(burstRequests);

	for req := range(burstRequests) { 
		t, ok := <-burstRateLimiter;
		if !ok  { 
			break
		}
		fmt.Println("request", req, t);
	}

	fmt.Println("Graceful shutdown complete")
}

func usuaulRateLimiter(wg *sync.WaitGroup) { 
	defer wg.Done();
	requests := make(chan int, 5);
	limiter := time.Tick(200 * time.Millisecond);
	for i:= range(5) { 
		requests <- i;
	}
	close(requests);

	for i := range(requests) { 
		<- limiter;
		fmt.Println("rate limiter", i, time.Now());
	}
}