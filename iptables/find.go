package iptables

import (
	"os"
	"os/exec"
)

var ip4tables string
var ip6tables string
var ip4tablesRestore string
var ip6tablesRestore string

const rulesV4Path = "/etc/iptables/rules.v4"
const rulesV6Path = "/etc/iptables/rules.v6"

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

func findIP4Tables() {
	if _, err := os.Stat(ip4tables); !os.IsNotExist(err) {
		return // then our path already exists, we have nothing to find
	}
	path, err := exec.LookPath("iptables")
	if err != nil {
		ip4tables = "iptables"
	}
	ip4tables = path
}

func findIP6Tables() {
	if _, err := os.Stat(ip6tables); !os.IsNotExist(err) {
		return // then our path already exists, we have nothing to find
	}
	path, err := exec.LookPath("ip6tables")
	if err != nil {
		ip6tables = "ip6tables"
	}
	ip6tables = path
}

func findIP4TablesRestore() {
	if _, err := os.Stat(ip4tablesRestore); !os.IsNotExist(err) {
		return // then our path already exists, we have nothing to find
	}
	path, err := exec.LookPath("iptables-restore")
	if err != nil {
		ip4tablesRestore = "iptables-restore"
	}
	ip4tablesRestore = path
}

func findIP6TablesRestore() {
	if _, err := os.Stat(ip6tablesRestore); !os.IsNotExist(err) {
		return // then our path already exists, we have nothing to find
	}
	path, err := exec.LookPath("ip6tables-restore")
	if err != nil {
		ip6tablesRestore = "ip6tables-restore"
	}
	ip6tablesRestore = path
}
