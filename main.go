package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aarnaud/http-mitigation/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/aarnaud/http-mitigation/server"
	"github.com/aarnaud/http-mitigation/db"
)

var Config *config.ServiceConfig
var cli = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		startCli()
	},
}

var cliOptionVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version.",
	Long:  "The version of this program",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version 0.0.1")
	},
}

func init() {
	cli.AddCommand(cliOptionVersion)

	flags := cli.Flags()

	flags.BoolP("verbose", "v", false, "Enable verbose")
	viper.BindPFlag("verbose", flags.Lookup("verbose"))

	flags.Int("listen-port", 8000, "HTTP listen port")
	viper.BindPFlag("listen_port", flags.Lookup("listen-port"))

	flags.String("cookie-name", "__mitigation", "Cookie Name")
	viper.BindPFlag("cookie_name", flags.Lookup("cookie-name"))

	flags.String("redis-addr", "127.0.0.1:6379", "Redis Server Address")
	viper.BindPFlag("redis_addr", flags.Lookup("redis-addr"))

	flags.String("redis-password", "", "Redis Password")
	viper.BindPFlag("redis_password", flags.Lookup("redis-password"))

	flags.Int("redis-db", 0, "Redis DB")
	viper.BindPFlag("redis_db", flags.Lookup("redis-db"))

	flags.Int("threshold1", 10000, "Threshold per domain per minute (mitigation redirect 307)")
	viper.BindPFlag("threshold1", flags.Lookup("threshold1"))

	flags.Int("threshold2", 50000, "Threshold per domain per minute (mitigation redirect javascript)")
	viper.BindPFlag("threshold2", flags.Lookup("threshold2"))
}

func main() {
	cli.Execute()
}

func startCli() {
	// EXPORT HM_LISTEN_PORT=8000
	viper.SetEnvPrefix("HM")
	viper.AutomaticEnv()

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.Info("Starting...")

	config.Config = &config.ServiceConfig{
		HTTPPort:   viper.GetInt("listen_port"),
		CookieName: viper.GetString("cookie_name"),
		RedisAddr: viper.GetString("redis_addr"),
		RedisPassword: viper.GetString("redis_password"),
		RedisDB: viper.GetInt("redis_db"),
		Threshold1: viper.GetInt64("threshold1"),
		Threshold2: viper.GetInt64("threshold2"),
	}

	log.Debugf("Config: %+v", config.Config)
	db.Connect()
	server.Start()
}