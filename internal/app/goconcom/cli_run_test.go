package goconcom

import (
	"github.com/roemer/gover"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

var examplesDir = filepath.Join("..", "..", "..", "examples")

func TestFindChangelogBottomUpWorkdirSameLevel(t *testing.T) {
	ass := assert.New(t)

	workingPath := filepath.Join(examplesDir, "subdir1")
	absolulteWorkingPath, err := filepath.Abs(workingPath)
	ass.Nil(err)

	changelogPath, err := findChangelogBottomUp(absolulteWorkingPath)

	ass.Nil(err)
	ass.Equal(filepath.Join(absolulteWorkingPath, ChangelogDefaultName), changelogPath)
}

func TestFindChangelogBottomUpWorkdirBelow(t *testing.T) {
	ass := assert.New(t)

	workingPath := filepath.Join(examplesDir, "subdir2", "subdir21")
	absolulteWorkingPath, err := filepath.Abs(workingPath)
	expectedChangelogPath, err := filepath.Abs(filepath.Join(examplesDir, "subdir2", ChangelogDefaultName))
	ass.Nil(err)

	changelogPath, err := findChangelogBottomUp(absolulteWorkingPath)

	ass.Nil(err)
	ass.Equal(expectedChangelogPath, changelogPath)
}

func TestGetVersionsFromChangelog_NoChangelog(t *testing.T) {
	ass := assert.New(t)

	_, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "UNKNOWN.md"))
	ass.Error(err)
}

func TestGetVersionsFromChangelog_EmtpyChangelog(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "EMTPY.md"))

	ass.Nil(err)
	ass.Empty(versions)
}

func TestGetVersionsFromChangelog_OnlyOnePatchVersion(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "ONLY_ONE_PATCH.md"))

	ass.Nil(err)
	ass.Equal(1, len(versions))
	ass.True(gover.ParseSimple(0, 0, 1).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_OnlyOneMinorVersion(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "ONLY_ONE_MINOR.md"))

	ass.Nil(err)
	ass.Equal(1, len(versions))
	ass.True(gover.ParseSimple(0, 1, 0).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_OnlyOneMajorVersion(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "ONLY_ONE_MAJOR.md"))

	ass.Nil(err)
	ass.Equal(1, len(versions))
	ass.True(gover.ParseSimple(1, 0, 0).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_SeveralPatchVersions(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "SEVERAL_PATCHES.md"))

	ass.Nil(err)
	ass.Equal(3, len(versions))
	ass.True(gover.ParseSimple(0, 0, 3).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_SeveralMinorVersions(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "SEVERAL_MINORS.md"))

	ass.Nil(err)
	ass.Equal(3, len(versions))
	ass.True(gover.ParseSimple(0, 3, 0).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_SeveralMajorVersions(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "SEVERAL_MAJORS.md"))

	ass.Nil(err)
	ass.Equal(3, len(versions))
	ass.True(gover.ParseSimple(3, 0, 0).Equals(versions[0]))
}

func TestGetVersionsFromChangelog_SeveralMixedVersions(t *testing.T) {
	ass := assert.New(t)

	versions, err := getVersionsFromChangelog(filepath.Join(examplesDir, "changelogs", "SEVERAL_MIXED.md"))

	ass.Nil(err)
	ass.Equal(9, len(versions))
	ass.True(gover.ParseSimple(1, 0, 1).Equals(versions[0]))
}
