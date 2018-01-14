package _sampleproject

import "github.com/sky0621/go-testcode-autogen/example/sampleproject/config"

type SmplApp struct {
	Cfg *config.Config
}

func (a *SmplApp) Start() (int, error) {
	return 1, nil
}
