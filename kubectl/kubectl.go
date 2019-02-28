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

func (k kubectl) GetPod(namespace string) string {
	params := []string{"--namespace", namespace, "get", "pod"}
	res, _ := k.runCmd(params, nil)
	return res
}

func (k kubectl) Apply(namespace string, content string) error {
	params := []string{"--namespace", namespace, "apply", "-f", "-"}
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