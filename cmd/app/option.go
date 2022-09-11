package app

import (
	"github.com/spf13/pflag"
)

//CliOptions abstractions configuration options for reading parameters from the
// command line.
type CliOptions interface {
	AddFlags(fs *pflag.FlagSet)
}

//ConfigurableOptions abstractions configuration options for reading parameters from
// the configuration file
type ConfigurableOptions interface {
	ApplyFlags() []error
}
