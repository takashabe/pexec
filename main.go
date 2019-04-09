package main

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/spf13/cobra"
)

var (
	command  string
	count    int
	parallel int
)

func main() {
	NewCommand().Execute()
}

func NewCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "pexec",
		Short: "Executions any command",
		Run:   execution,
	}

	root.PersistentFlags().StringVarP(&command, "exec", "e", "", "Exec command (required)")
	root.PersistentFlags().IntVarP(&count, "count", "c", 1, "Number of executions")
	root.PersistentFlags().IntVarP(&parallel, "parallel", "p", 1, "Number of parallel executions")
	root.MarkPersistentFlagRequired("exec")

	return root
}

func execution(cmd *cobra.Command, args []string) {
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

			ec := exec.Command("sh", "-c", command)
			ec.Stdout = os.Stdout
			ec.Stderr = os.Stderr
			if err := ec.Run(); err != nil {
				fmt.Printf("received error: %s, command: %s", err.Error(), command)
			}
		}()
	}
	wg.Wait()
}
