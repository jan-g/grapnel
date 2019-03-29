package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type Person struct {
	Name string `json:"name"`
	Sleep float64 `json:"sleep"`
	Dir string `json:"dir"`
	Cat string `json:"file"`
	Shell string `json:"shell"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := &Person{Name: "World"}
	json.NewDecoder(in).Decode(p)
	time.Sleep(time.Duration(p.Sleep * float64(time.Second)))
	e := os.Environ()
	msg := struct {
		Msg string `json:"message"`
		Env []string `json:"env"`
		Files []string `json:"files"`
		Err string `json:"error"`
		Content []byte `json:"bytes"`
		Stdout []byte `json:"stdout"`
		Stderr []byte `json:"stderr"`
	}{
		Msg: fmt.Sprintf("Hello %s", p.Name),
		Env: e,
	}
	if p.Dir != "" {
		files := []string{}
		if fs, err := ioutil.ReadDir(p.Dir); err == nil {
			for _, f := range fs {
				files = append(files, f.Name())
			}
			msg.Files = files
		} else {
			msg.Err = err.Error()
		}
	}
	if p.Cat != "" {
		if content, err := ioutil.ReadFile(p.Cat); err == nil {
			msg.Content = content
		} else {
			msg.Err = err.Error()
		}
	}
	if p.Shell != "" {
		cmd := exec.Command("/bin/sh", "-c", p.Shell)
		stdout, err1 := cmd.StdoutPipe()
		stderr, err2 := cmd.StderrPipe()
			outs := make(chan []byte, 1)
			errs := make(chan []byte, 1)
		if err1 != nil {
			msg.Err = err1.Error()
			goto done
		}
		if err2 != nil {
			msg.Err = err2.Error()
			goto done
		}

		go func() {
			result, _ := ioutil.ReadAll(stdout)
			stdout.Close()
			outs <- result
		}()

		go func() {
			result, _ := ioutil.ReadAll(stderr)
			stderr.Close()
			errs <- result
		}()

		if err := cmd.Run(); err != nil {
			msg.Err = err.Error()
			goto done
		}
		msg.Stdout = <- outs
		msg.Stderr = <- errs
		done:
	}
	json.NewEncoder(out).Encode(&msg)
}
