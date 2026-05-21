package integration_tests

import (
	gctxcmd "gctx/cmd"
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
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.NoActiveContextMessage+"\n"+"p1 p2\n")
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
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.CurrentContextMessage("p2")+"\n"+"p1 p2\n")
}

func TestChangesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "personal.config"), "[user]\n\tname = Bob\n\temail = bob@example.test\n")

	out, err := runGctx(t, repo, configDir, "personal")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.SwitchedContextMessage("personal")+"\n")
	requireEqual(t, strings.TrimSpace(runGit(t, repo, "config", "--local", "user.name")), "Bob")
	requireEqual(t, strings.TrimSpace(runGit(t, repo, "config", "--local", "gctx.profile")), "personal")
}

func TestChangingProfileReplacesActiveProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "p1.config"), "[user]\n\tname = One\n")
	writeFile(t, filepath.Join(configDir, "p2.config"), "[user]\n\tname = Two\n")

	out, err := runGctx(t, repo, configDir, "p1")
	requireNoError(t, err, out)
	out, err = runGctx(t, repo, configDir, "p2")
	requireNoError(t, err, out)

	got := strings.TrimSpace(runGit(t, repo, "config", "--local", "--get-all", "gctx.profile"))
	requireEqual(t, got, "p2")
}

func TestSavesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	runGit(t, repo, "config", "--local", "user.name", "Alice")
	runGit(t, repo, "config", "--local", "user.email", "alice@example.test")

	out, err := runGctx(t, repo, configDir, "work", "--save")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.SavedContextMessage("work")+"\n")
	content := string(mustReadFile(t, filepath.Join(configDir, "work.config")))
	requireContains(t, content, "name = Alice")
	requireContains(t, content, "email = alice@example.test")
}

func TestRemovesProfile(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "work.config"), "[user]\n\tname = Alice\n")

	out, err := runGctx(t, repo, configDir, "work", "--remove")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.RemovedContextMessage("work")+"\n")
	if _, err := os.Stat(filepath.Join(configDir, "work.config")); !os.IsNotExist(err) {
		t.Fatal("removed profile still exists", err)
	}
}

func TestRequiresProfileNameWithSaveOrRemove(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "--save")
	requireError(t, err, out)
	requireContains(t, out, gctxcmd.ProfileNameRequiredMessage())
}

func TestRejectsSaveAndRemoveTogether(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "work", "--save", "--remove")
	requireError(t, err, out)
	requireContains(t, out, gctxcmd.ConflictingActionMessage())
}

func TestSaveErrorsOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	configDir := t.TempDir()

	out, err := runGctx(t, dir, configDir, "work", "--save")
	requireError(t, err, out)
}

func TestChangeErrorsOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "work.config"), "[user]\n\tname = Alice\n")

	out, err := runGctx(t, dir, configDir, "work")
	requireError(t, err, out)
}

func TestChangeErrorsWhenProfileIsMissing(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	writeFile(t, filepath.Join(configDir, "p1.config"), "[user]\n\tname = One\n")
	writeFile(t, filepath.Join(configDir, "p2.config"), "[user]\n\tname = Two\n")

	out, err := runGctx(t, repo, configDir, "missing")
	requireError(t, err, out)
	requireContains(t, out, gctxcmd.MissingContextMessage("missing", []string{"p1", "p2"}))
}

func TestChangeErrorsWhenNoProfilesExist(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()

	out, err := runGctx(t, repo, configDir, "missing")
	requireError(t, err, out)
	requireContains(t, out, gctxcmd.MissingContextMessage("missing", nil))
}

func TestSaveErrorsWhenGitConfigIsMissing(t *testing.T) {
	repo := initGitRepo(t)
	configDir := t.TempDir()
	if err := os.Remove(filepath.Join(repo, ".git", "config")); err != nil {
		t.Fatal(err)
	}

	out, err := runGctx(t, repo, configDir, "work", "--save")
	requireError(t, err, out)
}

func TestGlobalListsProfilesWithActiveProfileOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	configDir := t.TempDir()
	writeFile(t, filepath.Join(home, ".gitconfig"), "[gctx]\n\tprofile = global-work\n")
	writeFile(t, filepath.Join(configDir, "global-work.config"), "[user]\n\tname = Global\n")

	out, err := runGctxWithHome(t, dir, home, configDir, "--global")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.CurrentContextMessage("global-work")+"\n"+"global-work\n")
}

func TestGlobalSavesProfileOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	configDir := t.TempDir()
	writeFile(t, filepath.Join(home, ".gitconfig"), "[user]\n\tname = Global Alice\n\temail = global@example.test\n")

	out, err := runGctxWithHome(t, dir, home, configDir, "global-work", "--save", "--global")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.SavedContextMessage("global-work")+"\n")
	content := string(mustReadFile(t, filepath.Join(configDir, "global-work.config")))
	requireContains(t, content, "name = Global Alice")
	requireContains(t, content, "email = global@example.test")
}

func TestGlobalChangesProfileOutsideGitRepo(t *testing.T) {
	dir := t.TempDir()
	home := t.TempDir()
	configDir := t.TempDir()
	writeFile(t, filepath.Join(home, ".gitconfig"), "[user]\n\tname = Old\n")
	writeFile(t, filepath.Join(configDir, "global-work.config"), "[user]\n\tname = Global Bob\n\temail = bob@example.test\n")

	out, err := runGctxWithHome(t, dir, home, configDir, "global-work", "-g")
	requireNoError(t, err, out)

	requireEqual(t, out, gctxcmd.SwitchedContextMessage("global-work")+"\n")
	globalConfig := string(mustReadFile(t, filepath.Join(home, ".gitconfig")))
	requireContains(t, globalConfig, "name = Global Bob")
	requireContains(t, globalConfig, "profile = global-work")
}

func runGctx(t *testing.T, dir, configDir string, args ...string) (string, error) {
	t.Helper()
	return runGctxWithHome(t, dir, "", configDir, args...)
}

func runGctxWithHome(t *testing.T, dir, home, configDir string, args ...string) (string, error) {
	t.Helper()
	allArgs := append([]string{"--config", configDir}, args...)
	cmd := exec.Command(gctxBin, allArgs...)
	cmd.Dir = dir
	if home != "" {
		cmd.Env = append(os.Environ(), "HOME="+home)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal("git failed", args, err, string(out))
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
		t.Fatal("missing text", want, "in", text)
	}
}

func requireNotContains(t *testing.T, text, unwanted string) {
	t.Helper()
	if strings.Contains(text, unwanted) {
		t.Fatal("unexpected text", unwanted, "in", text)
	}
}

func requireEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatal("got", got, "want", want)
	}
}

func requireNoError(t *testing.T, err error, out string) {
	t.Helper()
	if err != nil {
		t.Fatal(err, out)
	}
}

func requireError(t *testing.T, err error, out string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error", out)
	}
}
