// +build !prod

package kops


type MockHandler struct {
	Cmds chan string
	Responses chan string
	_    struct{}
}

func (k MockHandler) QueryCmd(paramString string, stdInData []byte) ([]byte, error) {
	k.Cmds <- paramString
	return []byte(<- k.Responses), nil
}

func (k MockHandler) RunCmd(paramString string, stdInData []byte) error {
	k.Cmds <- paramString
	return nil
}
