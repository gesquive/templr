package cmd

import (
	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:     "up",
	Aliases: []string{"load", "start"},
	Short:   "Bring up the firewall(s)",
	Long:    `Generates the firewall rules and activates them.`,
	Run:     runLoad,
}

func init() {
	RootCmd.AddCommand(loadCmd)

	// #viperbug
	// loadCmd.Flags().StringP("rules", "r", "",
	// 	"The templated firewall rules")

	// viper.BindEnv("rules")

	// viper.BindPFlag("rules", loadCmd.Flags().Lookup("rules"))
}

func runLoad(cmd *cobra.Command, args []string) {
	loadRules()
}
