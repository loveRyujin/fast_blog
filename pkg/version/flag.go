package version

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

type versionValue int

const (
	// 未设置版本
	VersionNotSet versionValue = 0
	// 启用版本
	VersionEnable versionValue = 1
	// 原始版本
	VersionRaw versionValue = 2
)

const strRawVersion string = "raw"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() any {
	return *v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	if boolVal {
		*v = VersionEnable
	} else {
		*v = VersionNotSet
	}
	return nil
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionEnable))
}

func (v *versionValue) Type() string {
	return "version"
}

func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	flag.Lookup(name).NoOptDefVal = "true"
}

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit")

func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

func PrintAndExitIfRequested() {
	// 检查版本信息的值并打印相应的信息
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)
	} else if *versionFlag == VersionEnable {
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)
	}
}
