package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	gctxBin string
	rootDir string
)

func TestMain(m *testing.M) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		os.Exit(1)
	}
	rootDir = filepath.Dir(filepath.Dir(file))

	tmpDir, err := os.MkdirTemp("", "gctx-tests-*")
	if err != nil {
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	gctxBin = filepath.Join(tmpDir, "gctx")
	cmd := exec.Command("go", "build", "-o", gctxBin, ".")
	cmd.Dir = rootDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestListsProfilesWithoutActiveProfileAndIgnoresTrash(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "p2.config"), "[user]\n\tname = Two\n")
	writeFile(t, filepath.Join(configDir, "trash.txt"), "ignore me")
	writeFile(t, filepath.Join(configDir, ".DS_Store"), "ignore me")
	if err := os.Mkdir(filepath.Join(configDir, "p3.config"), 0o700); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(configDir, "p1.config"), "[user]\n\tname = One\n")

	out, err := runGctx(t, repo, configDir)
	if err != nil {
		t.Fatalf("gctx failed: %v\n%s", err, out)
	}

	requireContains(t, out, "(didn't find active profile)")
	requireContains(t, out, "p1 p2")
	requireNotContains(t, out, "trash")
	requireNotContains(t, out, "p3")
}

func TestListsProfilesWithActiveProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	runGit(t, repo, "config", "--local", "gctx.profile", "p2")
	writeFile(t, filepath.Join(configDir, "p1.config"), "[user]\n\tname = One\n")
	writeFile(t, filepath.Join(configDir, "p2.config"), "[user]\n\tname = Two\n")

	out, err := runGctx(t, repo, configDir)
	if err != nil {
		t.Fatalf("gctx failed: %v\n%s", err, out)
	}

	requireContains(t, out, "current used profile:")
	requireContains(t, out, "p2")
	requireContains(t, out, "p1 p2")
}

func TestChangesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "personal.config"), "[user]\n\tname = Bob\n\temail = bob@example.test\n")

	out, err := runGctx(t, repo, configDir, "personal")
	if err != nil {
		t.Fatalf("gctx failed: %v\n%s", err, out)
	}

	requireContains(t, out, "personal changed successfully")
	if got := strings.TrimSpace(runGit(t, repo, "config", "--local", "user.name")); got != "Bob" {
		t.Fatalf("user.name = %q, want Bob", got)
	}
	if got := strings.TrimSpace(runGit(t, repo, "config", "--local", "gctx.profile")); got != "personal" {
		t.Fatalf("gctx.profile = %q, want personal", got)
	}
}

func TestChangingProfileReplacesActiveProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "p1.config"), "[user]\n\tname = One\n")
	writeFile(t, filepath.Join(configDir, "p2.config"), "[user]\n\tname = Two\n")

	if out, err := runGctx(t, repo, configDir, "p1"); err != nil {
		t.Fatalf("gctx p1 failed: %v\n%s", err, out)
	}
	if out, err := runGctx(t, repo, configDir, "p2"); err != nil {
		t.Fatalf("gctx p2 failed: %v\n%s", err, out)
	}

	got := strings.TrimSpace(runGit(t, repo, "config", "--local", "--get-all", "gctx.profile"))
	if got != "p2" {
		t.Fatalf("gctx.profile values = %q, want only p2", got)
	}
}

func TestSavesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	runGit(t, repo, "config", "--local", "user.name", "Alice")
	runGit(t, repo, "config", "--local", "user.email", "alice@example.test")

	out, err := runGctx(t, repo, configDir, "work", "--save")
	if err != nil {
		t.Fatalf("gctx failed: %v\n%s", err, out)
	}

	requireContains(t, out, "work saved successfully")
	content := string(mustReadFile(t, filepath.Join(configDir, "work.config")))
	requireContains(t, content, "name = Alice")
	requireContains(t, content, "email = alice@example.test")
}

func TestRemovesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "work.config"), "[user]\n\tname = Alice\n")

	out, err := runGctx(t, repo, configDir, "work", "--remove")
	if err != nil {
		t.Fatalf("gctx failed: %v\n%s", err, out)
	}

	requireContains(t, out, "work removed successfully")
	if _, err := os.Stat(filepath.Join(configDir, "work.config")); !os.IsNotExist(err) {
		t.Fatalf("removed profile still exists, stat err: %v", err)
	}
}

func TestRequiresProfileNameWithSaveOrRemove(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "--save")
	if err == nil {
		t.Fatalf("gctx --save succeeded, want error\n%s", out)
	}
	requireContains(t, out, "profile name is required")
}

func TestRejectsSaveAndRemoveTogether(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "work", "--save", "--remove")
	if err == nil {
		t.Fatalf("gctx --save --remove succeeded, want error\n%s", out)
	}
	requireContains(t, out, "use either --save or --remove")
}

func TestSaveErrorsOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	configDir := t.TempDir()

	out, err := runGctx(t, dir, configDir, "work", "--save")
	if err == nil {
		t.Fatalf("gctx save outside git repo succeeded, want error\n%s", out)
	}
	requireContains(t, out, "couldn't find the git repo")
}

func TestChangeErrorsOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "work.config"), "[user]\n\tname = Alice\n")

	out, err := runGctx(t, dir, configDir, "work")
	if err == nil {
		t.Fatalf("gctx change outside git repo succeeded, want error\n%s", out)
	}
	requireContains(t, out, "couldn't find the git repo")
}

func TestChangeErrorsWhenProfileIsMissing(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "missing")
	if err == nil {
		t.Fatalf("gctx missing profile succeeded, want error\n%s", out)
	}
	requireContains(t, out, "failed to read profile missing")
}

func TestSaveErrorsWhenGitConfigIsMissing(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	if err := os.Remove(filepath.Join(repo, ".git", "config")); err != nil {
		t.Fatal(err)
	}

	out, err := runGctx(t, repo, configDir, "work", "--save")
	if err == nil {
		t.Fatalf("gctx save with missing git config succeeded, want error\n%s", out)
	}
	requireContains(t, out, "could not read file")
}

func runGctx(t *testing.T, dir, configDir string, args ...string) (string, error) {
	t.Helper()

	allArgs := append([]string{"--config", configDir}, args...)
	cmd := exec.Command(gctxBin, allArgs...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
	return string(out)
}

func initGitRepo(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	runGit(t, dir, "init")
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func requireContains(t *testing.T, text, want string) {
	t.Helper()

	if !strings.Contains(text, want) {
		t.Fatalf("%q does not contain %q", text, want)
	}
}

func requireNotContains(t *testing.T, text, unwanted string) {
	t.Helper()

	if strings.Contains(text, unwanted) {
		t.Fatalf("%q contains %q", text, unwanted)
	}
}
