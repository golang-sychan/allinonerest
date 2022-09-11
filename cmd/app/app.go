package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
)

var (
	progressMessage = color.BlueString("==>")
	usageTemplate   = fmt.Sprintf(`%s{{if .Runnable}}
  %s{{end}}{{if .HasAvailableSubCommands}}
  %s{{end}}{{if gt (len .Aliases) 0}}

%s
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

%s{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%s --help" for more information about a command.{{end}}
`,
		color.CyanString("Usage:"),
		color.GreenString("{{.UseLine}}"),
		color.GreenString("{{.CommandPath}} [command]"),
		color.CyanString("Aliases:"),
		color.CyanString("Examples:"),
		color.CyanString("Available Commands:"),
		color.GreenString("{{rpad .Name .NamePadding }}"),
		color.CyanString("Flags:"),
		color.CyanString("Global Flags:"),
		color.CyanString("Additional help topics:"),
		color.GreenString("{{.CommandPath}} [command]"),
	)
)

// App application management
type App struct {
	basename    string
	name        string
	description string
	logger      interface{}
	options     CliOptions
	commands    []*Command
	runfunc     RunFunc
	//rootcmd is the root cobra command for this app
	rootcmd *cobra.Command
}

// Option represent functions apply the specific option to the app
//, as the same function signature that we can define any kind behavior,
// such as func(app *App) { app.basename = "foo" }, func(app *App) { app.logger = zrus.logger }
type Option func(*App)

// WithDescription is used to set the description of the application.
func WithDescription(desc string) Option {
	return func(a *App) {
		a.description = desc
	}
}

// WithOptions to open the application's function to read from the command line
// or read parameters from the configuration file.
func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

type RunFunc func(basename string) error

// WithRunFunc construct a func to return func(*App) as type Option, hence it can be used in
// construct App as options in parameters
func WithRunFunc(run RunFunc) Option {
	return func(app *App) {
		app.runfunc = run
	}
}

func NewApp(name, basename string, opts ...Option) *App {
	a := &App{
		basename: basename,
		name:     name,
	}
	//apply the configurations
	for _, o := range opts {
		o(a)
	}
	//
	return a
}

//intiRootCmd 初始化处理 cobra root command
func (a *App) intiRootCmd() {
	initFlag()

	cmd := cobra.Command{
		Use:           a.basename,
		Long:          a.description,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetUsageTemplate(usageTemplate)
	cmd.SetOut(os.Stdout)
	cmd.Flags().SortFlags = false
	for _, command := range a.commands {
		cmd.AddCommand(command.cobraCommand())
	}
	cmd.SetHelpCommand(helpCommand(a.name))
	if a.runfunc != nil {
		cmd.Run = a.runCommand
	}
	if a.options != nil {
		if _, ok := a.options.(ConfigurableOptions); ok {
			addConfigFlag(a.basename, cmd.Flags())
		}
		a.options.AddFlags(cmd.Flags())
	}
	addHelpFlag(a.name, cmd.Flags())

	a.rootcmd = &cmd
}

func (a *App) runCommand(cmd *cobra.Command, args []string) {
	// merge configuration and print it
	if a.options != nil {
		if configurableOptions, ok := a.options.(ConfigurableOptions); ok {
			if errs := configurableOptions.ApplyFlags(); len(errs) > 0 {
				for _, err := range errs {
					fmt.Printf("%v %v\n", color.RedString("Error:"), err)
				}
				os.Exit(1)
			}
			printConfig()
		}
	}
}
