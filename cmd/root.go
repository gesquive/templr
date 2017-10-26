package cmd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gesquive/cli"
	"github.com/gesquive/shield/engine"
	"github.com/gesquive/shield/iptables"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var cfgFile string
var displayVersion string

var verbose bool
var logDebug bool
var showVersion bool

var runIPv4 bool
var runIPv6 bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "shield",
	Short:            "iptables firewall manager",
	Long:             `Manage and update your iptables firewall rules`,
	PersistentPreRun: preRun,
	Hidden:           true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/%s\n",
		RootCmd.HelpTemplate(), displayVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		"config file (default is $HOME/.config/shield.yml)")
	RootCmd.PersistentFlags().StringP("log-file", "l", "",
		"Path to log file")

	RootCmd.PersistentFlags().BoolP("ipv4-only", "4", false,
		"Apply command to IPv4 rules only.")
	RootCmd.PersistentFlags().BoolP("ipv6-only", "6", false,
		"Apply command to IPv6 rules only.")

	// This is a workaround for https://github.com/spf13/viper/issues/233
	//TODO: remove this once bug is fixed #viperbug
	RootCmd.PersistentFlags().StringP("rules", "r", "",
		"The templated firewall rules")

	RootCmd.PersistentFlags().BoolVarP(&logDebug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("shield")
	viper.AutomaticEnv()

	viper.BindEnv("ipv4-only")
	viper.BindEnv("ipv6-only")
	viper.BindEnv("rules")

	viper.BindPFlag("ipv4-only", RootCmd.PersistentFlags().Lookup("ipv4-only"))
	viper.BindPFlag("ipv6-only", RootCmd.PersistentFlags().Lookup("ipv6-only"))
	viper.BindPFlag("rules", RootCmd.PersistentFlags().Lookup("rules"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		homeConfig := path.Join(home, ".config/shield")

		viper.SetConfigName("config")   // name of config file (without extension)
		viper.AddConfigPath(".")        // adding current directory as first search path
		viper.AddConfigPath(homeConfig) // adding home directory as next search path
		viper.AddConfigPath("/etc/shield")

	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func preRun(cmd *cobra.Command, args []string) {
	if logDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	logFilePath := getLogFilePath(viper.GetString("log_file"))
	log.Debugf("config: log_file=%s", logFilePath)
	if len(logFilePath) == 0 {
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: true,
		})
		log.SetOutput(os.Stdout)
	} else {
		log.SetFormatter(&prefixed.TextFormatter{
			TimestampFormat: time.RFC3339,
		})

		logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file=%v", err)
		}
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	log.Debugf("config: file=%s", viper.ConfigFileUsed())

	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
	}
	cli.Debug("Running with debug turned on")

	ipv4Only := viper.GetBool("ipv4-only")
	ipv6Only := viper.GetBool("ipv6-only")
	if ipv4Only == ipv6Only {
		runIPv4 = true
		runIPv6 = true
	} else if ipv4Only {
		runIPv4 = true
		runIPv6 = false
	} else if ipv6Only {
		runIPv4 = false
		runIPv6 = true
	}
	//TODO: check if we have permissions to change iptables
}

func getLogFilePath(defaultPath string) (logPath string) {
	fi, err := os.Stat(defaultPath)
	if err == nil && fi.IsDir() {
		logPath = path.Join(defaultPath, "shield.log")
	} else {
		logPath = defaultPath
	}
	return
}

func loadRules() {
	rulePath := viper.GetString("rules")
	fmt.Printf("rules=%s\n", rulePath)
	if len(rulePath) == 0 {
		cli.Error("No rules specified")
		os.Exit(2)
	}

	//TODO: Check that DNS works
	rules, err := engine.NewRuleset(rulePath)
	if err != nil {
		log.Error("%v", err)
		os.Exit(2)
	}

	b, err := rules.GenerateRules(displayVersion, []string{})
	if err != nil {
		log.Error("%v", err)
		os.Exit(2)
	}

	if runIPv4 {
		cli.Info("Applying IPv4 firewall rules")
		err := iptables.LoadIPv4Rules(b)
		if err != nil {
			log.Error("%v", err)
			os.Exit(10)
		}
	}

	if runIPv6 {
		cli.Info("Applying IPv6 firewall rules")
		err := iptables.LoadIPv6Rules(b)
		if err != nil {
			log.Error("%v", err)
			os.Exit(10)
		}
	}
}

func unloadRules() {
	if runIPv4 {
		iptables.ClearIPv4Rules()
	}

	if runIPv6 {
		iptables.ClearIPv6Rules()
	}
}
