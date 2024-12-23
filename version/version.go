package version

var (
	// Version is the version of this software.
	// This value is overwritten by the build script.
	Version = "0.0.0"

	// CommitHash is the git commit hash of this software.
	// This value is overwritten by the build script.
	CommitHash = "unknown"

	// ShortCommitHash is the short version of the commit hash.
	ShortCommitHash = "unknown"
)

func init() {
	if len(CommitHash) < 7 {
		ShortCommitHash = CommitHash
	} else {
		ShortCommitHash = CommitHash[:7]
	}
}
