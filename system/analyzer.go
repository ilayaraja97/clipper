package system

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

const APPLICATION_NAME = "Clipper"

type Analysis struct {
	operatingSystem OperatingSystem
	distribution    string
	shell           string
	homeDirectory   string
	username        string
	editor          string
	configFile      string
}

func (a *Analysis) GetApplicationName() string {
	return APPLICATION_NAME
}

func (a *Analysis) GetOperatingSystem() OperatingSystem {
	return a.operatingSystem
}

func (a *Analysis) GetDistribution() string {
	return a.distribution
}

func (a *Analysis) GetShell() string {
	return a.shell
}

func (a *Analysis) GetHomeDirectory() string {
	return a.homeDirectory
}

func (a *Analysis) GetUsername() string {
	return a.username
}

func (a *Analysis) GetEditor() string {
	return a.editor
}

func (a *Analysis) GetConfigFile() string {
	return a.configFile
}

func Analyse() *Analysis {
	return &Analysis{
		operatingSystem: GetOperatingSystem(),
		distribution:    GetDistribution(),
		shell:           GetShell(),
		homeDirectory:   GetHomeDirectory(),
		username:        GetUsername(),
		editor:          GetEditor(),
		configFile:      GetConfigFile(),
	}
}

func GetOperatingSystem() OperatingSystem {
	switch runtime.GOOS {
	case "linux":
		return LinuxOperatingSystem
	case "darwin":
		return MacOperatingSystem
	case "windows":
		return WindowsOperatingSystem
	default:
		return UnknownOperatingSystem
	}
}

func GetDistribution() string {
	if runtime.GOOS != "linux" {
		return ""
	}

	content, err := os.ReadFile(filepath.Clean("/etc/os-release"))
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(content), "\n") {
		if !strings.HasPrefix(line, "PRETTY_NAME=") {
			continue
		}

		dist := strings.TrimPrefix(line, "PRETTY_NAME=")
		return strings.Trim(dist, "\"")
	}

	return ""
}

func GetShell() string {
	shell := strings.TrimSpace(os.Getenv("SHELL"))
	if shell != "" {
		split := strings.Split(shell, "/")

		return split[len(split)-1]
	}

	if runtime.GOOS == "windows" {
		return "powershell"
	}

	return "sh"
}

func GetHomeDirectory() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		return ""
	}

	return homeDir
}

func GetUsername() string {
	name := strings.TrimSpace(os.Getenv("USER"))
	if name == "" {
		currentUser, err := user.Current()
		if err == nil && currentUser.Username != "" {
			if split := strings.Split(currentUser.Username, `\`); len(split) > 0 {
				name = split[len(split)-1]
			} else {
				name = currentUser.Username
			}
		}
	}

	return strings.TrimSpace(name)
}

func GetEditor() string {
	name := strings.TrimSpace(os.Getenv("EDITOR"))
	if name != "" {
		return strings.TrimSpace(name)
	}

	if runtime.GOOS == "windows" {
		return "notepad"
	}

	return "nano"
}

func GetConfigFile() string {
	return filepath.Join(
		GetHomeDirectory(),
		".config",
		strings.ToLower(APPLICATION_NAME)+".json",
	)
}
