package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	gitVersion   = "v0.0.0-master+$FORMAT:%H$"
	gitCommit    = "$FORMAT:%H$"
	gitTreeState = ""

	buildDate = "1970-01-01T00:00:00Z"
)

// Info contains version information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

func (i Info) String() string {
	return i.GitVersion
}

// ToJson以json格式输出版本信息
func (i Info) ToJson() string {
	s, _ := json.Marshal(i)
	return string(s)
}

// ToText以文本格式输出版本信息
func (i Info) ToText() string {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("GitVersion:", i.GitVersion)
	table.AddRow("GitCommit:", i.GitCommit)
	table.AddRow("GitTreeState:", i.GitTreeState)
	table.AddRow("BuildDate:", i.BuildDate)
	table.AddRow("GoVersion:", i.GoVersion)
	table.AddRow("Compiler:", i.Compiler)
	table.AddRow("Platform:", i.Platform)

	return table.String()
}

// Get返回详细的版本信息
func Get() Info {
	return Info{
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
