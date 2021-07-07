package vm

import (
	"bytes"
	"os/exec"
)

func VboxCommandHandler(params ...string) (string, error) {
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command("vboxmanage", params...)
	cmd.Stderr = &errBuf
	cmd.Stdout = &outBuf
	if err := cmd.Run(); err != nil {
		return outBuf.String(), err
	}
	return outBuf.String(), nil
}
