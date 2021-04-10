package main

import (
	"os/user"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/Depado/scarecrow/cmd"
)

// Main command that will be run when no other command is provided on the
// command-line
var rootCmd = &cobra.Command{
	Use: "scarecrow",
	Run: func(c *cobra.Command, args []string) {
		fx.New(
			fx.NopLogger,
			fx.Provide(
				cmd.NewConf, cmd.NewLogger,
			),
			fx.Invoke(hw),
		).Run()
	},
}

func hw(c *cmd.Conf, l zerolog.Logger) {
	l.Info().Msg("Hello world")
	usr, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to get user dir")
	}
	log.Info().Str("home", usr.HomeDir).Msg("User home found")
}

func main() {
	// Initialize Cobra and Viper
	cmd.AddAllFlags(rootCmd)
	rootCmd.AddCommand(cmd.VersionCmd)

	// Run the command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("unable to start")
	}
}
