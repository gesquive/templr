package iptables

import (
	"os"

	sh "github.com/codeskyblue/go-sh"
)

func getCleanupRules() []byte {
	return []byte(`
*filter
:INPUT ACCEPT [0:0]
:FORWARD ACCEPT [0:0]
:OUTPUT ACCEPT [0:0]
COMMIT
`)
}

func Exists() bool {
	if _, err := os.Stat(ip4tables); !os.IsNotExist(err) {
		return true
	}
	return false
}

func LoadIPv4Rules(rules []byte, persist bool) error {
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

	if persist {
		err = persistIPv4Rules(rules)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadIPv6Rules(rules []byte, persist bool) error {
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

	if persist {
		err = persistIPv6Rules(rules)
		if err != nil {
			return err
		}
	}

	return nil
}

func persistIPv4Rules(rules []byte) error {
	err := writeFile(ip4RulesPersistPath, rules)
	if err != nil {
		return err
	}
	return nil
}

func persistIPv6Rules(rules []byte) error {
	err := writeFile(ip6RulesPersistPath, rules)
	if err != nil {
		return err
	}
	return nil
}

func ClearIPv4Rules(persist bool) error {
	// First set all chains to accept in case something funky happens
	cleanRules := getCleanupRules()
	err := LoadIPv4Rules(cleanRules, persist)
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

func ClearIPv6Rules(persist bool) error {
	// First set all chains to accept in case something funky happens
	cleanRules := getCleanupRules()
	err := LoadIPv6Rules(cleanRules, persist)
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
	out, err := sh.Command(ip4tables, "-L", "-v", "-n").Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func GetIPv6Summary() string {
	out, err := sh.Command(ip6tables, "-L", "-v", "-n").Output()
	if err != nil {
		return ""
	}
	return string(out)
}
