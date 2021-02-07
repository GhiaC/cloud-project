package vm

import (
	"bytes"
	"github.com/ghiac/go-commons/log"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func VboxCommandHandler(params ...string) (string, error) {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("vboxmanage", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = &outBuf
	if err := cmd.Run(); err != nil {
		log.Logger.
			WithFields(logrus.Fields{
				"location": "VboxCommandHandler",
				"params":   params,
			}).
			Error(errBuf.String())
		return outBuf.String(), err
	}
	return outBuf.String(), nil
}
