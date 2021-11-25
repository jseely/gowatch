package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var reqPerDuration = flag.Duration("r", time.Second, "# of requests per second, must be time.ParseDuration parsable")

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Printf("usage: %s [-r <duration string>] <executable> [<arg>]+", os.Args[0])
		return
	}
	fmt.Printf("Running `%s` every %s\n", strings.Join(flag.Args(), " "), reqPerDuration.String())

	output := make(chan string, 10)
	done := make(chan bool)

	go func() {
		i := 0
		for {
			s, more := <-output
			fmt.Printf("%d\t: %s\n", i, s)
			i++
			if !more {
				done <- true
				break
			}
		}
	}()

	for {
		go func(out chan string) {
			outBytes, err := exec.Command(flag.Args()[0], flag.Args()[1:]...).Output()
			if err != nil {
				out <- err.Error()
				return
			}
			out <- strings.TrimSpace(string(outBytes))
		}(output)
		time.Sleep(*reqPerDuration)
	}
}
