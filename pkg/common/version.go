package common

import "fmt"

var (
	NAME    = "stackit-nuke"
	SUMMARY = "0.0.0-dev"
	BRANCH  = "dev"
	COMMIT  = "dirty"
)

type AppVersionInfo struct {
	Name    string
	Branch  string
	Summary string
	Commit  string
}

func (a *AppVersionInfo) String() string {
	return fmt.Sprintf("%s - %s - %s", a.Name, a.Summary, a.Commit)
}

var AppVersion AppVersionInfo

func init() {
	AppVersion = AppVersionInfo{
		Name:    NAME,
		Branch:  BRANCH,
		Summary: SUMMARY,
		Commit:  COMMIT,
	}
}
