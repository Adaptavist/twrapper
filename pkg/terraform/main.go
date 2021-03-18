package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
)

func writeJSONFile(path string, data interface{}) (err error) {
	content, err := json.Marshal(data)

	if err != nil {
		return
	}

	err = ioutil.WriteFile(path, content, 0644)

	return
}

// createBackendConfigFile creates file configuring your backend
func createBackendConfigFile(opts Opts) error {
	return writeJSONFile("terraform.backend.tf.json", Root{
		Terraform: Terraform{
			Backend: opts.Backend,
		},
	})
}

// createVariablesFile creates a variables file for the Terraform run
func createVariablesFile(opts Opts) error {
	return writeJSONFile("terraform.tfvars.json", opts.Variables)
}

func Configure(opts Opts) (err error) {
	if !opts.Backend.IsEmpty() {
		fmt.Println("writing backend")
		if err = createBackendConfigFile(opts); err != nil {
			return
		}
	}

	err = createVariablesFile(opts)
	return
}

func copyLogs(r io.Reader) {
	buf := make([]byte, 80)
	for {
		n, err := r.Read(buf)
		if n > 0 {
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

func Init() error {
	return run("terraform", "init")
}

func Execute(opts Opts) error {
	return run("terraform", opts.Arguments...)
}

// NewBackend creates a structure suitable for marshalling into a terraform file
func NewBackend(backendType string, backendConfig map[string]interface{}) Backend {
	if len(backendConfig) == 0 {
		backendConfig = make(map[string]interface{})
	}
	return Backend{
		backendType: backendConfig,
	}
}
