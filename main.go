package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

var (
	command  string
	count    int
	parallel int
)

func main() {
	flag.StringVar(&command, "e", "", "require: the exec command.")
	flag.IntVar(&count, "c", 1, "optional: number of executions. default 1.")
	flag.IntVar(&parallel, "p", 1, "optional: number of parallel executions. default 1(not parallels).")
	flag.Parse()

	semaphore := make(chan struct{}, parallel)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			semaphore <- struct{}{}
			defer func() {
				<-semaphore
				wg.Done()
			}()

			cmd := exec.Command("sh", "-c", command)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("received error: %s, command: %s", err.Error(), command)
			}
		}()
	}
	wg.Wait()
}
