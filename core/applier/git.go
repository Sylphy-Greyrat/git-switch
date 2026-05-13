package applier

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sylphy/git-switch/core/config"
)

type GitApplier struct {
	configPath string
}

func NewGitApplier(configPath string) GitApplier {
	return GitApplier{configPath: configPath}
}

func (a GitApplier) ApplyGitConfig(ctx context.Context, profile config.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := config.ExpandHome(a.configPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create git config directory: %w", err)
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read git config %s: %w", path, err)
	}

	merged := mergeGitConfig(string(existing), profile)

	if err := os.WriteFile(path, []byte(merged), 0o644); err != nil {
		return fmt.Errorf("write git config %s: %w", path, err)
	}
	return nil
}

func mergeGitConfig(existing string, profile config.Profile) string {
	lines := strings.Split(strings.TrimSuffix(existing, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}

	userStart, userEnd := findSection(lines, "user", "")
	lines = setSectionKey(lines, userStart, userEnd, "user", "", "name", profile.User.Name)
	lines = setSectionKey(lines, userStart, userEnd, "user", "", "email", profile.User.Email)

	if profile.SSH != nil && profile.SSH.HostAlias != "" && len(profile.SSH.Hosts) > 0 {
		subname := fmt.Sprintf("git@%s:", profile.SSH.HostAlias)
		urlStart, urlEnd := findSection(lines, "url", subname)
		insteadOf := fmt.Sprintf("git@%s:", profile.SSH.Hosts[0])
		lines = setSectionKey(lines, urlStart, urlEnd, "url", subname, "insteadOf", insteadOf)
	}

	return strings.Join(lines, "\n") + "\n"
}

func findSection(lines []string, section, subname string) (start, end int) {
	header := "[" + section
	if subname != "" {
		header += fmt.Sprintf(" %q", subname)
	}
	header += "]"

	for i, line := range lines {
		if strings.TrimSpace(line) == header {
			start = i + 1
			end = start
			for end < len(lines) {
				trimmed := strings.TrimSpace(lines[end])
				if strings.HasPrefix(trimmed, "[") {
					break
				}
				end++
			}
			return start, end
		}
	}
	return -1, -1
}

func setSectionKey(lines []string, start, end int, section, subname, key, value string) []string {
	newLine := "\t" + key + " = " + value
	prefix := key + " ="

	if start >= 0 {
		// Section exists — find and replace the key
		for i := start; i < end; i++ {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), prefix) {
				lines[i] = newLine
				return lines
			}
		}
		// Key not found — append to existing section before next section or at end
		lines = append(lines, "")
		copy(lines[end+1:], lines[end:])
		lines[end] = newLine
		return lines
	}

	// Section doesn't exist — append
	header := "[" + section
	if subname != "" {
		header += fmt.Sprintf(" %q", subname)
	}
	header += "]"

	if len(lines) > 0 && lines[len(lines)-1] != "" {
		lines = append(lines, "")
	}
	lines = append(lines, header, newLine)
	return lines
}
