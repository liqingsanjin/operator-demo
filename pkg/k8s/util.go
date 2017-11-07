package k8s

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	ServiceAccountNSFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

// GetCurrentNS gets the namepace of the current Pod.
// It could be used in a Pod only.
func GetCurrentNS() (string, error) {
	// If setting `automountServiceAccountToken` to `false` for PodSpec, then
	// `ServiceAccountNSFile` will not exist and reading a non-exist file will panic.
	if _, err := os.Stat(ServiceAccountNSFile); err != nil && os.IsNotExist(err) {
		return "", errors.New("couldn't get the current namespace")
	}

	dataBytes, err := ioutil.ReadFile(ServiceAccountNSFile)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(dataBytes)), nil
}