package goconcom

import (
	"flag"
	"fmt"
	"github.com/roemer/goconcom/pkg/logging"
	"github.com/roemer/gover"
	"github.com/samber/lo"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
)

const ChangelogDefaultName = "CHANGELOG.md"

func RunCmd(args []string) error {
	// Flags and help for the command
	var verbose bool
	var workingDirectory string

	flagSetRun := flag.NewFlagSet("run", flag.ExitOnError)
	flagSetRun.BoolVar(&verbose, "verbose", false, "The flag to set in order to get verbose output.")
	flagSetRun.BoolVar(&verbose, "v", verbose, "Alias for -verbose.")
	flagSetRun.StringVar(&workingDirectory, "workDir", ".", "The path to the working directory.")
	flagSetRun.Usage = func() { printCmdUsage(flagSetRun, "run", "") }
	flagSetRun.Parse(args)

	// Create a logger
	desiredLogLevel := lo.Ternary(verbose, slog.LevelDebug, slog.LevelInfo)
	logger := slog.New(logging.NewReadableTextHandler(os.Stdout, &logging.ReadableTextHandlerOptions{Level: desiredLogLevel}))
	logger.Debug(fmt.Sprintf("Initialized logger with level: %s", desiredLogLevel))
	logger.Info("Starting gonovate run")

	// Traverse the working directory and find the nearest CHANGELOG.md
	logger.Info(fmt.Sprintf("Working directory: %s", workingDirectory))
	changelogPath, err := findChangelogBottomUp(workingDirectory)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Found changelog on path: %s", changelogPath))

	logger.Info("Extracting version informations from the changelog")
	versions, err := getVersionsFromChangelog(changelogPath)
	if err != nil {
		return err
	}

	highestVersion := gover.FindMax(versions, gover.EmptyVersion, true)
	logger.Info(fmt.Sprintf("Found highest version: %d.%d.%d", highestVersion.Major(), highestVersion.Minor(), highestVersion.Patch()))

	return nil
}

func findChangelogBottomUp(startDir string) (string, error) {
	currentDir := startDir

	for {
		changelogPath := filepath.Join(currentDir, "CHANGELOG.md")
		if _, err := os.Stat(changelogPath); err == nil {
			return changelogPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}
	return "", fmt.Errorf("CHANGELOG.md not found")
}

func getVersionsFromChangelog(path string) ([]*gover.Version, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	md := goldmark.New()

	// Parse the Markdown content
	reader := text.NewReader(content)
	doc := md.Parser().Parse(reader)

	headings := []string{}

	// Traverse the AST to find headings
	if err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if heading, ok := n.(*ast.Heading); ok && entering {
			text := heading.Text(content)
			headings = append(headings, string(text))
			fmt.Printf("Heading: %s\n", text)
		}
		return ast.WalkContinue, nil
	}); err != nil {
		return nil, err
	}

	// Filter headings and convert them to versions
	versionRegex, err := regexp.Compile(`v?(\d+)(?:\.(\d+))?(?:\.(\d+))?`)
	if err != nil {
		return nil, err
	}
	versions := []*gover.Version{}
	for _, heading := range headings {
		// pre-validate match before version parsing?
		if versionRegex.MatchString(heading) {
			version, err := gover.ParseVersionFromRegex(heading, versionRegex)
			if err != nil {
				return nil, err
			}
			versions = append(versions, version)
		}
	}

	return versions, nil
}
