package iptables

import (
	"os"

	sh "github.com/codeskyblue/go-sh"
)

func Exists() bool {
	if _, err := os.Stat(ip4tables); !os.IsNotExist(err) {
		return true
	}
	return false
}

func LoadIPv4Rules(rules []byte) error {
	rulesFile, err := getTempFile()
	if err != nil {
		return err
	}
	defer os.Remove(rulesFile.Name())

	err = writeFile(rulesFile.Name(), rules)
	if err != nil {
		return err
	}

	err = sh.Command(ip4tablesRestore, rulesFile.Name()).Run()
	if err != nil {
		return err
	}

	return nil
}

func LoadIPv6Rules(rules []byte) error {
	rulesFile, err := getTempFile()
	if err != nil {
		return err
	}
	defer os.Remove(rulesFile.Name())

	err = writeFile(rulesFile.Name(), rules)
	if err != nil {
		return err
	}

	err = sh.Command(ip6tablesRestore, rulesFile.Name()).Run()
	if err != nil {
		return err
	}

	return nil
}

func ClearIPv4Rules() error {
	// First set all chains to accept in case something funky happens
	cleanRules := []byte(`
*filter
:INPUT ACCEPT [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
COMMIT
`)

	err := LoadIPv4Rules(cleanRules)
	if err != nil {
		return err
	}

	// flush the nat & mangle tables
	sh.Command(ip4tables, "-t", "nat", "-F").Run()
	sh.Command(ip4tables, "-t", "mangle", "-F").Run()
	// flush all chains
	sh.Command(ip4tables, "-F").Run()
	// delete all non-default chains
	sh.Command(ip4tables, "-X").Run()

	return nil
}

func ClearIPv6Rules() error {
	// First set all chains to accept in case something funky happens
	cleanRules := []byte(`
*filter
:INPUT ACCEPT [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
COMMIT
`)

	err := LoadIPv6Rules(cleanRules)
	if err != nil {
		return err
	}

	// flush the nat & mangle tables
	sh.Command(ip6tables, "-t", "nat", "-F").Run()
	sh.Command(ip6tables, "-t", "mangle", "-F").Run()
	// flush all chains
	sh.Command(ip6tables, "-F").Run()
	// delete all non-default chains
	sh.Command(ip6tables, "-X").Run()

	return nil
}

func GetIPv4Summary() string {
	out, err := sh.Command(ip4tables, "-L", "-n").Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func GetIPv6Summary() string {
	out, err := sh.Command(ip6tables, "-L", "-n").Output()
	if err != nil {
		return ""
	}
	return string(out)
}
