package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func shouldPrint() bool {
	return os.Getenv("HIDE_TF_OUTPUT") == ""
}

func writeJSONFile(path string, data interface{}) (err error) {
	content, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return
	}

	err = os.WriteFile(path, content, 0644)
	log.Printf("writing %s:\n%s\n", path, content)
	return
}

func copyLogs(r io.Reader) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 && shouldPrint() {
			fmt.Print(string(buf[0:n]))
		}
		if err != nil {
			break
		}
	}
}

func run(command string, args ...string) (err error) {
	fmt.Printf("%s %s\n", command, strings.Join(args, " "))
	var wg sync.WaitGroup

	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		copyLogs(stdout)
	}()

	go func() {
		defer wg.Done()
		copyLogs(stderr)
	}()

	wg.Wait()

	if err = cmd.Wait(); err != nil {
		return err
	}

	return
}
