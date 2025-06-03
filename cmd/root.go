/*
 */
package cmd

import (
	"context"
	"os"
	"time"

	"github.com/senzing-garage/demo-entity-search/httpserver"
	"github.com/senzing-garage/go-cmdhelping/cmdhelper"
	"github.com/senzing-garage/go-cmdhelping/option"
	"github.com/senzing-garage/go-cmdhelping/option/optiontype"
	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Long string = `
demo-entity-search long description.
    `
	ReadHeaderTimeout        = 60
	Short             string = "demo-entity-search short description"
	Use               string = "demo-entity-search"
)

var avoidServe = option.ContextVariable{
	Arg:     "avoid-serving",
	Default: option.OsLookupEnvBool("SENZING_TOOLS_AVOID_SERVING", false),
	Envar:   "SENZING_TOOLS_AVOID_SERVING",
	Help:    "Avoid serving.  For testing only. [%s]",
	Type:    optiontype.Bool,
}

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	avoidServe,
	option.EnableAll,
	option.HTTPPort,
	option.ServerAddress,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command.
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command.
func RunE(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	httpServer := &httpserver.BasicHTTPServer{
		AvoidServing:      viper.GetBool(avoidServe.Arg),
		EnableAll:         viper.GetBool(option.EnableAll.Arg),
		ReadHeaderTimeout: ReadHeaderTimeout * time.Second,
		ServerAddress:     viper.GetString(option.ServerAddress.Arg),
		ServerPort:        viper.GetInt(option.HTTPPort.Arg),
	}

	err := httpServer.Serve(ctx)

	return wraperror.Errorf(err, wraperror.NoMessage)
}

// Used in construction of cobra.Command.
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}
