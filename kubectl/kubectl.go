package kubectl

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os/exec"
)

type kubectl struct {
	context string
	_ struct{}
}

func New(context string) kubectl {
	return kubectl{context:context}
}

func (k kubectl) GetPods(namespace string) string {
	var ns []string
	if namespace != "" {
		ns = []string{"--namespace", namespace}
	}
	params := append(ns, "get", "pod")
	res, _ := k.runCmd(params, nil)
	return res
}

func (k kubectl) Apply(namespace string, content string) error {
	var ns []string
	if namespace != "" {
		ns = []string{"--namespace", namespace}
	}
	params := append(ns, "apply", "-f", "-")
	_, err := k.runCmd(params, []byte(content))
	return err
}

func (k kubectl) runCmd(params []string, stdInData []byte)  (string, error) {
	ctx := []string{"--context", k.context}
	cmd := exec.Command("kubectl", append(ctx, params...)...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if stdInData != nil {
		cmd.Stdin = bytes.NewBuffer(stdInData)
	}

	_ = cmd.Start()

	if errb.String() != "" {
		log.Fatalln(errb.String())
	}
	e := cmd.Wait()

	return outb.String(), e
}

func (k kubectl) printOut(out io.ReadCloser) {
	scanner := bufio.NewScanner(out)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		log.Println(scanner.Text())
	}
}