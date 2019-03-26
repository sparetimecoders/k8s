// +build !prod

package kops

type MockHandler struct {
	Cmds chan string
	_    struct{}
}

func (k MockHandler) QueryCmd(paramString string, stdInData []byte) ([]byte, error) {
	k.Cmds <- paramString
	return []byte{}, nil
}

func (k MockHandler) RunCmd(paramString string, stdInData []byte) error {
	k.Cmds <- paramString
	return nil
}

func (k MockHandler) MinimumKopsVersionInstalled(requiredKopsVersion string) bool {
	return true
}
