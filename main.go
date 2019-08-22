package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	MaxConcurrentThread = 8
)

func init() {
	if i, err := strconv.Atoi(os.Getenv("MAX_CONCURRENT_THREAD")); err == nil && 0 < i {
		MaxConcurrentThread = i
	}
}

func main() {
	var stdin string
	w := csv.NewWriter(os.Stdout)

	var mu sync.Mutex
	write := func(input []string) {
		mu.Lock()
		defer mu.Unlock()
		w.Write(input)
		w.Flush()
	}

	write([]string{
		"host",
		"error",
		"expires_on",
	})

	var (
		wg = sync.WaitGroup{}
		ch = make(chan int, MaxConcurrentThread)
	)
	for {
		n, err := fmt.Scan(&stdin)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		ch <- 1
		wg.Add(1)
		go func(input string) {
			defer func() {
				<-ch
				wg.Done()
			}()
			r, err := SSLCheck(stdin)
			if err != nil {
				log.Printf("[error]%s --> %v\n", stdin, err)
				write([]string{
					stdin,
					fmt.Sprint(err),
					"",
				})
			} else {
				write([]string{
					fmt.Sprintf("%s:%s", r.Host, r.Port),
					fmt.Sprint(r.Error),
					fmt.Sprint(r.ExpiresOn),
				})
			}
		}(stdin)
	}
	wg.Wait()
}
