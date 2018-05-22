package cmd

import (
	"fmt"
	"net"
	"os"
	"path"
	"time"

	"github.com/gesquive/cli"
	"github.com/gesquive/templr/engine"
	"github.com/gesquive/templr/iptables"
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
var persist bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "templr",
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
		"config file (default is $HOME/.config/templr.yml)")
	RootCmd.PersistentFlags().StringP("log-file", "l", "",
		"Path to log file")

	RootCmd.PersistentFlags().BoolP("ipv4-only", "4", false,
		"Apply command to IPv4 rules only.")
	RootCmd.PersistentFlags().BoolP("ipv6-only", "6", false,
		"Apply command to IPv6 rules only.")
	RootCmd.PersistentFlags().BoolP("persist", "p", false,
		"Save the firewall configuration to netfilter-persistent")

	// This is a workaround for https://github.com/spf13/viper/issues/233
	//TODO: remove this once bug is fixed #viperbug
	RootCmd.PersistentFlags().StringP("rules", "r", "",
		"The templated firewall rules")

	RootCmd.PersistentFlags().BoolVarP(&logDebug, "debug", "D", false,
		"Write debug messages to console")
	RootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false,
		"Show the version and exit")

	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("templr")
	viper.AutomaticEnv()

	viper.BindEnv("ipv4-only")
	viper.BindEnv("ipv6-only")
	viper.BindEnv("persist")
	viper.BindEnv("rules")

	viper.BindPFlag("ipv4-only", RootCmd.PersistentFlags().Lookup("ipv4-only"))
	viper.BindPFlag("ipv6-only", RootCmd.PersistentFlags().Lookup("ipv6-only"))
	viper.BindPFlag("persist", RootCmd.PersistentFlags().Lookup("persist"))
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
		homeConfig := path.Join(home, ".config/templr")

		viper.SetConfigName("config")   // name of config file (without extension)
		viper.AddConfigPath(".")        // adding current directory as first search path
		viper.AddConfigPath(homeConfig) // adding home directory as next search path
		viper.AddConfigPath("/etc/templr")

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
		log.SetOutput(logFile)
	}
	log.Debugf("config: log_file=%s", logFilePath)
	log.Debugf("config: file=%s", viper.ConfigFileUsed())

	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
	}
	log.Debug("Running with debug turned on")

	persist = viper.GetBool("persist")
	log.Debugf("config: persist=%t", persist)

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
	log.Debugf("config: runIPv4=%t runIPv6=%t", runIPv4, runIPv6)

	if err := iptables.FindIPv4(); runIPv4 && err != nil {
		cli.Error("%s", err)
		cli.Error("Make sure iptables is installed")
		os.Exit(6)
	}
	if err := iptables.FindIPv6(); runIPv6 && err != nil {
		cli.Error("%s", err)
		cli.Error("Make sure ip6tables is installed")
		os.Exit(6)
	}

	if !isRootUser() {
		cli.Error("Modifying iptables requires root access")
		os.Exit(5)
	}
}

func getLogFilePath(defaultPath string) (logPath string) {
	fi, err := os.Stat(defaultPath)
	if err == nil && fi.IsDir() {
		logPath = path.Join(defaultPath, "templr.log")
	} else {
		logPath = defaultPath
	}
	return
}

func isRootUser() bool {
	uid := os.Geteuid()
	return uid == 0
}

func isDNSWorking() bool {
	addrs, err := net.LookupHost("github.com")
	if err != nil {
		return false
	}
	return len(addrs) != 0
}

func loadRules() {
	rulePath := viper.GetString("rules")
	if len(rulePath) == 0 {
		cli.Error("No rules specified")
		os.Exit(2)
	}

	if !isDNSWorking() {
		cli.Error("DNS is not resolving")
		os.Exit(4)
	}

	rules, err := engine.NewRuleset(rulePath)
	if err != nil {
		log.Error("%v", err)
		os.Exit(2)
	}

	data, err := rules.GenerateRules(displayVersion)
	if err != nil {
		log.Error("%v", err)
		os.Exit(2)
	}

	// right now, don't see a reason to make this an option
	restoreCounters := true

	if runIPv4 {
		log.Info("Applying IPv4 firewall rules")
		err := iptables.LoadIPv4Rules(data, restoreCounters, persist)
		if err != nil {
			log.Error("%v", err)
			os.Exit(10)
		}
	}

	if runIPv6 {
		log.Info("Applying IPv6 firewall rules")
		err := iptables.LoadIPv6Rules(data, restoreCounters, persist)
		if err != nil {
			log.Error("%v", err)
			os.Exit(10)
		}
	}
}

func unloadRules() {
	if runIPv4 {
		iptables.ClearIPv4Rules(persist)
	}

	if runIPv6 {
		iptables.ClearIPv6Rules(persist)
	}
}
