package cmd

import "github.com/spf13/cobra"

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:     "reload",
	Aliases: []string{"restart", "update"},
	Short:   "Reload the firewall rules",
	Long:    `Clear, regenerate and load the firewall rules`,
	Run:     runReload,
}

func init() {
	RootCmd.AddCommand(reloadCmd)

	// #viperbug
	// reloadCmd.Flags().StringP("rules", "r", "",
	// 	"The templated firewall rules")
	//
	// viper.BindEnv("rules")
	//
	// viper.BindPFlag("rules", reloadCmd.Flags().Lookup("rules"))
}

func runReload(cmd *cobra.Command, args []string) {
	unloadRules()
	loadRules()
}
