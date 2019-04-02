package kubectl

import (
	"bufio"
	"bytes"
	"log"
	"os/exec"
	"strings"
)

type KubeCtl interface {
	Apply(yamlContent string) error
}

type CmdHandler interface {
	ApplyCmd(data []byte) ([]byte, error)
}

type osCmdHandler struct {
	context string
	_       struct{}
}

type kubectl struct {
	Handler CmdHandler
	_       struct{}
}

func New(clusterName string) KubeCtl {
	return kubectl{Handler: osCmdHandler{context: clusterName}}
}

func (c kubectl) Apply(yamlContent string) error {
	parts := getParts(yamlContent)
	log.Printf("Found %d parts in yaml content", len(parts))
	for _, part := range parts {
		if res, err := c.Handler.ApplyCmd([]byte(part)); err != nil {
			log.Printf("Failed to apply content:\n %v\n\nReason: %v, %v\n", part, err, string(res))
			return err
		}
	}
	return nil
}

func (h osCmdHandler) ApplyCmd(data []byte) ([]byte, error) {
	cmd := exec.Command("kubectl", "--context", h.context, "apply", "-f", "-")
	if data != nil {
		cmd.Stdin = bytes.NewBuffer(data)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

// getParts returns the different YAML parts in a string
// The triple-dash '---' is used as the separator
func getParts(yamlContent string) []string {
	scanner := bufio.NewScanner(strings.NewReader(yamlContent))
	var yamls, current []string
	for scanner.Scan() {
		if line := scanner.Text(); line == "---" {
			yamls = append(yamls, strings.Join(current, "\n"))
			current = nil
		} else {
			current = append(current, line)
		}
	}
	if len(current) > 0 {
		yamls = append(yamls, strings.Join(current, "\n"))
	}
	return yamls
}
