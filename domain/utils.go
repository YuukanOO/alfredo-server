package domain

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/go-playground/validator.v9"
)

// Validator for structs of the domain
var validate = validator.New()

func computeCheckSum(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}

func isComponent(str string) bool {
	return str[:1] == "<"
}

const transformScript = "console.log(require('babel-core').transform(process.argv[1], { plugins: ['transform-react-jsx', 'transform-es2015-arrow-functions'],}).code);"

func getWidgetBytes(relativeDir string, entry string) ([]byte, error) {
	if isComponent(entry) {
		return []byte(entry), nil
	}

	data, err := ioutil.ReadFile(filepath.Join(relativeDir, entry))

	if err != nil {
		return nil, err
	}

	return data, nil
}

func transformWidget(widget string) (string, error) {
	cmd := exec.Command("node", "-e", transformScript, widget)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		return stderr.String(), err
	}

	return fmt.Sprintf("function(device, command, showView) { return %s; }", strings.TrimSpace(stdout.String())), nil
}
