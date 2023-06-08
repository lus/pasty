package meta

import "strings"

const devEnvironmentName = "dev"

var (
	Environment = devEnvironmentName
	Version     = "dev"
)

func IsProdEnvironment() bool {
	return strings.ToLower(Environment) != devEnvironmentName
}
