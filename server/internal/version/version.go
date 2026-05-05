package version

var (
	Branch  = "unknown"
	Commit  = "unknown"
	Tag     = ""
	RepoURL = "https://github.com/ai-hermes/easyssl"
)

func String() string {
	if Tag != "" {
		return Tag
	}
	shortCommit := Commit
	if len(shortCommit) > 7 {
		shortCommit = shortCommit[:7]
	}
	return shortCommit
}

func DetailString() string {
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

func BranchURL() string {
	if RepoURL == "" || Branch == "" || Branch == "unknown" {
		return ""
	}
	return RepoURL + "/tree/" + Branch
}

func TagURL() string {
	if RepoURL == "" || Tag == "" {
		return ""
	}
	return RepoURL + "/tree/" + Tag
}
