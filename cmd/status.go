package cmd

import (
	"github.com/gesquive/cli"
	"github.com/gesquive/shield/iptables"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Report the firewall status",
	Long:  `Retrieves a summary of the current firewall(s).`,
	Run:   runStatus,
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	if runIPv4 {
		cli.Info("IPv4 Firewall Status")
		cli.Info("----------------------------------------------------")
		status := iptables.GetIPv4Summary()
		cli.Info(status)
	}
	if runIPv6 {
		cli.Info("IPv6 Firewall Status")
		cli.Info("----------------------------------------------------")
		status := iptables.GetIPv6Summary()
		cli.Info(status)
	}
}
