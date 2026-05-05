package version

var (
	Branch = "unknown"
	Commit = "unknown"
)

func String() string {
	return Branch + " @ " + Commit
}
