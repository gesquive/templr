package cmd

import "github.com/spf13/cobra"

// unloadCmd represents the unload command
var unloadCmd = &cobra.Command{
	Use:     "unload",
	Aliases: []string{"down", "stop", "clear"},
	Short:   "Clear the firewall, accept all traffic",
	Long:    `Clear the firewall, accept all traffic.`,
	Run:     runUnload,
}

func init() {
	RootCmd.AddCommand(unloadCmd)
}

func runUnload(cmd *cobra.Command, args []string) {
	unloadRules()
}
