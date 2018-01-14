package cmd

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/sky0621/go-testcode-autogen/example/sampleproject"
	"github.com/sky0621/go-testcode-autogen/example/sampleproject/config"
)

func main() {
	code := realMain()
	logrus.Infof("ExitCode: %v", code)
	os.Exit(int(code))
}

func realMain() (exitCode int) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Errorf("Panic occured. ERR: %+v", err)
			exitCode = 9
		}
	}()

	return wrappedMain()
}

func wrappedMain() int {
	err := config.ReadConfig("config/path.tml")
	if err != nil {
		logrus.Errorf("[wrappedMain][call config.ReadConfig()] %#v\n", err)
		return 3
	}

	app := sampleproject.SmplApp{Cfg: config.NewConfig()}
	res, err := app.Start()
	if err != nil {
		logrus.Errorf("[wrappedMain][call app.Start()] %#v\n", err)
		return 4
	}
	return res
}
