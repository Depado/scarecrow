package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddLoggerFlags adds support to configure the level of the logger.
func AddLoggerFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.type", "console", `one of "console" or "json"`)
	c.PersistentFlags().Bool("log.caller", false, "display the file and line where the call was made")
}

// AddConfigurationFlag adds support to provide a configuration file on the
// command line.
func AddConfigurationFlag(c *cobra.Command) {
	c.PersistentFlags().StringP("conf", "c", "", "configuration file to use")
}

// AddChiaFlags adds support for various chia related command line flags
func AddChiaFlags(c *cobra.Command) {
	c.PersistentFlags().String("chia.network", "mainnet", "chia network configuration to use")
}

// AddAllFlags will add all the flags provided in this package to the provided
// command and will bind those flags with viper.
func AddAllFlags(c *cobra.Command) {
	AddConfigurationFlag(c)
	AddLoggerFlags(c)
	AddChiaFlags(c)

	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		log.Fatal().Err(err).Msg("unable to bind flags")
	}
}
