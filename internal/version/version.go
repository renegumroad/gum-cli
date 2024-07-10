package version

var (
	VERSION = "development"
	RELEASE = "false"
)

func IsRelease() bool {
	return RELEASE == "true"
}
