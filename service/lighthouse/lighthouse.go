// Package lighthouse provides the lighthouse auditing service.
package lighthouse

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/giantswarm/microerror"
)

// AuditURL creates a lighthouse report and returns the path
func AuditURL(url, name, formFactor string, dockerLinks []string, ignoreCertErrors bool) (path string, err error) {
	fmt.Printf("Creating lighthouse report\nURL: %s\nForm factor: %s\nOutput file: %s.json\n", url, formFactor, name)

	pwd, err := os.Getwd()
	if err != nil {
		return "", microerror.Mask(err)
	}

	tmpDir, err := ioutil.TempDir("/tmp", "lighthouse-temp")
	if err != nil {
		return "", microerror.Mask(err)
	}
	defer os.RemoveAll(tmpDir)

	if formFactor != "desktop" && formFactor != "mobile" {
		formFactor = "desktop"
	}

	ignoreCertErrorsFlag := ""
	if ignoreCertErrors {
		ignoreCertErrorsFlag = "--ignore-certificate-errors"
	}

	linkArgs := []string{}
	for _, l := range dockerLinks {
		linkArgs = append(linkArgs, fmt.Sprintf("--link=%s", l))
	}

	args := []string{
		"run",
		"--rm",
		"--tty",
		fmt.Sprintf("-v=%s:/workdir", pwd),
		fmt.Sprintf("-v=%s:/dev/shm", tmpDir),
		"-w=/workdir",
	}

	moreArgs := []string{
		"quay.io/giantswarm/lighthouse:latest",
		"lighthouse",
		"--quiet",
		"--no-enable-error-reporting",
		"--output=json",
		fmt.Sprintf("--chrome-flags=--no-sandbox --headless %s", ignoreCertErrorsFlag),
		fmt.Sprintf("--emulated-form-factor=%s", formFactor),
		fmt.Sprintf("--output-path=/workdir/%s.json", name),
		url,
	}

	for _, a := range linkArgs {
		args = append(args, a)
	}
	for _, a := range moreArgs {
		args = append(args, a)
	}

	command := exec.Command("docker", args...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	if err != nil {
		_, errStr := string(stdout.Bytes()), string(stderr.Bytes())
		fmt.Printf("%s\n", errStr)
		return "", microerror.Mask(fmt.Errorf("cmd.Run() failed with %s", err))
	}

	return fmt.Sprintf("%s.json", name), nil
}
