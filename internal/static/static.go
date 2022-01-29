package static

// These variables represent the values that may be changed using ldflags
var (
	Version                   = "dev"
	EnvironmentVariablePrefix = "PASTY_"

	// TempFrontendPath defines the path where pasty loads the web frontend from; it will be removed any time soon
	// TODO: Remove this when issue #37 is fixed
	TempFrontendPath = "./web"
)
