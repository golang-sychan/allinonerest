package app

import (
	"github.com/spf13/pflag"
	"strings"
)

func initFlag() {
	pflag.CommandLine.SetNormalizeFunc(wordSeparatorNormalizeFunc)
}

func wordSeparatorNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}
