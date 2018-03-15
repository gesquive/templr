package iptables

import (
	"os/exec"

	"github.com/pkg/errors"
)

var ip4tables string
var ip6tables string
var ip4tablesRestore string
var ip6tablesRestore string

const ip4RulesPersistPath = "/etc/iptables/rules.v4"
const ip6RulesPersistPath = "/etc/iptables/rules.v6"

func SetIP4TablesPath(path string) {
	ip4tables = path
}

func SetIP6TablesPath(path string) {
	ip6tables = path
}

func SetIP4TablesRestorePath(path string) {
	ip4tablesRestore = path
}

func SetIP6TablesRestorePath(path string) {
	ip6tablesRestore = path
}

func Find() error {
	if err := FindIPv4(); err != nil {
		return err
	}
	if err := FindIPv6(); err != nil {
		return err
	}
	return nil
}

func FindIPv4() error {
	var err error
	if ip4tables, err = findUsableExe("iptables"); err != nil {
		return err
	}
	if ip4tablesRestore, err = findUsableExe("iptables-restore"); err != nil {
		return err
	}
	return nil
}

func FindIPv6() error {
	var err error
	if ip6tables, err = findUsableExe("ip6tables"); err != nil {
		return err
	}
	if ip6tablesRestore, err = findUsableExe("ip6tables-restore"); err != nil {
		return err
	}
	return nil
}

func findUsableExe(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}
	if !isExeUsable(path) {
		return "", errors.Errorf("Path is not executable %s", path)
	}
	return path, nil
}

func isExeUsable(path string) bool {
	err := exec.Command("test", "-x", path).Run()
	if err != nil {
		return false
	}
	return true
}
