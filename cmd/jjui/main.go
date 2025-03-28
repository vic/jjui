package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/idursun/jjui/internal/config"
	"github.com/idursun/jjui/internal/ui/context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/idursun/jjui/internal/ui"
)

var Version = "unknown"

var (
	versionFlag = flag.Bool("version", false, "show version")
	configFlag  = flag.Bool("config", false, "edit config")
)

func getVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
		return info.Main.Version
	}
	return Version
}

func main() {
	flag.Parse()
	if *versionFlag {
		println(getVersion())
		os.Exit(0)
	}
	if *configFlag {
		exitCode := config.Edit()
		os.Exit(exitCode)
	}

	var location string
	if len(os.Args) > 1 {
		location = os.Args[1]
	} else {
		location = os.Getenv("PWD")
	}

	rootLocation, err := getJJRootDir(location)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: There is no jj repo in \"%s\".\n", location)
		os.Exit(1)
	}

	appContext := context.NewAppContext(rootLocation)

	p := tea.NewProgram(ui.New(appContext), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func getJJRootDir(location string) (string, error) {
	cmd := exec.Command("jj", "root")
	cmd.Dir = location
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
