package version

var (
	Branch  = "unknown"
	Commit  = "unknown"
	RepoURL = ""
)

func String() string {
	shortCommit := Commit
	if len(shortCommit) > 7 {
		shortCommit = shortCommit[:7]
	}
	return Branch + " @ " + shortCommit
}

func CommitURL() string {
	if RepoURL == "" || Commit == "" || Commit == "unknown" {
		return ""
	}
	return RepoURL + "/commit/" + Commit
}
