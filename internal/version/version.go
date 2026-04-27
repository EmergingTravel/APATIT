package version

// These constants define application metadata.
// They can be overridden at build time using ldflags.
// Example: go build -ldflags "-X 'apatit/internal/version.Version=1.1.0'"
var (
	Name    = "apatit"
	Version = "v1.1.3"
	Owner   = "emerging.travel"
)

const (
	Language = "go"
)
