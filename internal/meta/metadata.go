package meta

import "strings"

const devEnvironmentName = "dev"

var (
	Environment = devEnvironmentName
)

func IsProdEnvironment() bool {
	return strings.ToLower(Environment) != devEnvironmentName
}
