package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompletionFilePathBash(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := completionFilePath("bash")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".local", "share", "git-switch", "completions", "git-switch.bash")
	if path != want {
		t.Fatalf("got %q, want %q", path, want)
	}
}

func TestCompletionFilePathZsh(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path, err := completionFilePath("zsh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(path, "_git-switch") {
		t.Fatalf("got %q, want path ending with _git-switch", path)
	}
}

func TestCompletionFilePathPwsh(t *testing.T) {
	path, err := completionFilePath("pwsh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(path, "git-switch.ps1") {
		t.Fatalf("got %q, want path ending with git-switch.ps1", path)
	}
}

func TestCompletionFilePathUnsupportedShell(t *testing.T) {
	_, err := completionFilePath("fish")
	if err == nil {
		t.Fatal("expected error for unsupported shell")
	}
}

func TestWriteAndRemoveCompletionScript(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := WriteCompletionScript("bash", "#!/bin/bash\necho test"); err != nil {
		t.Fatal(err)
	}
	path, _ := completionFilePath("bash")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "#!/bin/bash\necho test" {
		t.Fatalf("got %q", string(data))
	}

	if err := RemoveCompletionScript("bash"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("completion script should be deleted")
	}
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("empty completion dir should be deleted")
	}
}

func TestInjectCompletionBlockBash(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	rc := "export PATH=$PATH:/usr/local/bin\n"
	result, err := InjectCompletionBlock(rc, "bash")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, completionBlockBegin) {
		t.Fatal("should contain completion block begin")
	}
	if !strings.Contains(result, completionBlockEnd) {
		t.Fatal("should contain completion block end")
	}
	if !strings.Contains(result, "source ") {
		t.Fatal("should contain source line")
	}
}

func TestInjectCompletionBlockZsh(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	rc := "export PATH=$PATH:/usr/local/bin\n"
	result, err := InjectCompletionBlock(rc, "zsh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "fpath=(") {
		t.Fatal("should contain fpath")
	}
	// zsh block should NOT contain compinit (handled by user's own zshrc)
	if strings.Contains(result, "compinit") {
		t.Fatal("should not contain compinit")
	}
}

func TestRemoveCompletionBlock(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	rc := "export PATH=$PATH:/usr/local/bin\n# git-switch completion BEGIN\nsource /tmp/test.sh\n# git-switch completion END\necho done\n"
	result := RemoveCompletionBlock(rc, "bash")
	if strings.Contains(result, "git-switch completion") {
		t.Fatal("should not contain completion markers")
	}
	if !strings.Contains(result, "export PATH") {
		t.Fatal("should preserve original content")
	}
	if !strings.Contains(result, "echo done") {
		t.Fatal("should preserve trailing content")
	}
}

func TestIsCompletionInstalled(t *testing.T) {
	if IsCompletionInstalled("\n# git-switch completion BEGIN\nsource test.sh\n# git-switch completion END") != true {
		t.Fatal("should detect installed completion")
	}
	if IsCompletionInstalled("export PATH=/usr/bin") != false {
		t.Fatal("should return false when not installed")
	}
}

func TestInjectCompletionBlockIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	rc := "# git-switch completion BEGIN\nsource /tmp/test.sh\n# git-switch completion END\n"
	result, err := InjectCompletionBlock(rc, "bash")
	if err != nil {
		t.Fatal(err)
	}
	count := strings.Count(result, completionBlockBegin)
	if count != 1 {
		t.Fatalf("should be idempotent, got %d occurrences", count)
	}
}
