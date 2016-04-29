package commands

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

/*
 * Executes the given command in the given directory.
 *
 * To be sure that the NPM command does not encounter any problem, the
 * command output is checked to check that it does not contain "ERROR".
 *
 * TODO Use another solution to avoid loading all the output in memory.
 */
func ExecuteNPMCommand(cmd *exec.Cmd, workDir string) error {
	cmd.Dir = workDir

	var b bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &b)
	cmd.Stderr = io.MultiWriter(os.Stderr, &b)
	err := cmd.Run()

	if err != nil {
		return nil
	}

	/*
	 * Checking the command's output to detect errors.
	 */
	if strings.Contains(b.String(), "ERROR") {
		return errors.New("Build failed.")
	}

	return nil
}
