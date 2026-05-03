package command

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestByVersion(t *testing.T) {
	versions := []string{"1.2.3", "0.1.0", "1.10.2", "1.1.9"}
	expected := []string{"0.1.0", "1.1.9", "1.2.3", "1.10.2"}

	sort.Sort(ByVersion(versions))

	for i, v := range versions {
		if v != expected[i] {
			t.Errorf("at index %d: expected %s, got %s", i, expected[i], v)
		}
	}
}

func TestByVersionPanic(t *testing.T) {
	// First element invalid
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for invalid first version")
			}
		}()
		bv := ByVersion{"not-a-version", "1.0.0"}
		bv.Less(0, 1)
	}()

	// Second element invalid
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic for invalid second version")
			}
		}()
		bv := ByVersion{"1.0.0", "not-a-version"}
		bv.Less(0, 1)
	}()
}

func TestSetWellKnownValues(t *testing.T) {
	dict := make(map[string]string)
	dict = SetWellKnownValues(dict)

	keys := []string{"latesttag", "actualdate", "actualtime", "actualtimestamp"}
	for _, k := range keys {
		if _, ok := dict[k]; !ok {
			t.Errorf("expected key %s to be set", k)
		}
	}

	// Verify date format
	today := time.Now().Format("2006-01-02")
	if dict["actualdate"] != today {
		t.Errorf("expected date %s, got %s", today, dict["actualdate"])
	}
}

func TestGetGitLatestTag(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "versifyr-git-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	oldWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(oldWd) }()

	// 1. Test when not a git repo
	_, err = GetGitLatestTag()
	if err == nil {
		t.Error("expected error when not in a git repo")
	}

	// 2. Initialize a git repo
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Test when no tags
	tag, err := GetGitLatestTag()
	if err != nil {
		t.Errorf("unexpected error when no tags: %v", err)
	}
	if tag != "unknown" {
		t.Errorf("expected unknown when no tags, got %s", tag)
	}

	// 4. Add some commits and tags
	w, _ := repo.Worktree()
	_, _ = os.Create("file.txt")
	_, _ = w.Add("file.txt")
	_, _ = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})

	tags := []string{"1.0.0", "1.1.0", "2.0.0", "v3.0.0", "not-a-version"}
	for _, tg := range tags {
		head, _ := repo.Head()
		_, _ = repo.CreateTag(tg, head.Hash(), &git.CreateTagOptions{
			Tagger:  &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
			Message: tg,
		})
	}

	tag, err = GetGitLatestTag()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if tag != "2.0.0" {
		t.Errorf("expected 2.0.0, got %s", tag)
	}
}
