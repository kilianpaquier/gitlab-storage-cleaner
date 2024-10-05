package cobra

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// log is logrus default private logger retrieved into a variable.
	//
	// There's no need to give it to artifacts functions (v1, v2) Run since it's a shared pointer between logrus private and this variable.
	//
	// Unless later on there's a need to abstract Run loggers with other loggers
	// but it shouldn't since it's a CLI and it's an internal package.
	log = logrus.StandardLogger()

	logLevel  = "info"
	logFormat = "text"

	rootCmd = &cobra.Command{
		Use:               "gitlab-storage-cleaner",
		SilenceErrors:     true, // errors are already logged by fatal function when Execute has an error
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error { return preRun() },
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set logging level")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "text", `set logging format (either "text" or "json")`)

	_ = preRun() // ensure logging is correctly configured with default values even when a bad input flag is given
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fatal(context.Background(), err)
	}
}

func preRun() error {
	switch logFormat {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{
			DisableLevelTruncation: true,
			ForceColors:            true,
			FullTimestamp:          true,
			TimestampFormat:        time.RFC3339,
		})
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		return errors.New(`invalid --log-format argument, must be either "json" or "text"`)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	return nil
}

func fatal(ctx context.Context, err error) {
	log.WithContext(ctx).Fatal(err)
}
