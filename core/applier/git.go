package applier

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/sylphy/git-switch/core/config"
)

type GitApplier struct {
	configPath string
	backup     []byte
	applied    bool
}

var ErrApplyAlreadyActive = errors.New("git config already applied; revert before applying again")

func NewGitApplier(configPath string) GitApplier {
	return GitApplier{configPath: configPath}
}

func (a *GitApplier) ApplyGitConfig(ctx context.Context, profile config.Profile) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if a.applied {
		return ErrApplyAlreadyActive
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

	// Save backup for revert
	a.backup = existing

	merged := mergeGitConfig(string(existing), profile)

	if err := os.WriteFile(path, []byte(merged), 0o644); err != nil {
		return fmt.Errorf("write git config %s: %w", path, err)
	}
	a.applied = true
	return nil
}

func (a *GitApplier) Revert() error {
	if !a.applied {
		return nil
	}
	path, err := config.ExpandHome(a.configPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, a.backup, 0o644); err != nil {
		return fmt.Errorf("restore git config backup: %w", err)
	}
	a.backup = nil
	a.applied = false
	return nil
}

func mergeGitConfig(existing string, profile config.Profile) string {
	lines := strings.Split(strings.TrimSuffix(existing, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}

	// Remove duplicate [user] sections — only keep the first
	lines = removeDuplicateSections(lines, "user", "")

	// Find section after dedup
	userStart, userEnd := findSection(lines, "user", "")
	lines = setSectionKey(lines, userStart, userEnd, "user", "", "name", profile.User.Name)
	lines = setSectionKey(lines, userStart, userEnd, "user", "", "email", profile.User.Email)

	// GPG signing
	if profile.GPG != nil && profile.GPG.SigningKey != "" {
		lines = setSectionKey(lines, userStart, userEnd, "user", "", "signingkey", profile.GPG.SigningKey)
		commitStart, commitEnd := findSection(lines, "commit", "")
		gpgSign := "false"
		if profile.GPG.SignCommits {
			gpgSign = "true"
		}
		lines = setSectionKey(lines, commitStart, commitEnd, "commit", "", "gpgsign", gpgSign)
	}

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

func removeDuplicateSections(lines []string, section, subname string) []string {
	header := "[" + section
	if subname != "" {
		header += fmt.Sprintf(" %q", subname)
	}
	header += "]"

	var duplicates []struct{ start, end int }
	first := -1
	for i := 0; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == header {
			if first == -1 {
				first = i
			} else {
				// Find end of this duplicate section
				end := i + 1
				for end < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[end]), "[") {
					end++
				}
				duplicates = append(duplicates, struct{ start, end int }{i, end})
			}
		}
	}
	// Remove duplicates in reverse order to avoid index shifting
	for i := len(duplicates) - 1; i >= 0; i-- {
		lines = append(lines[:duplicates[i].start], lines[duplicates[i].end:]...)
	}
	return lines
}

func setSectionKey(lines []string, start, end int, section, subname, key, value string) []string {
	newLine := "\t" + key + " = " + value
	prefix := key + " ="

	if start >= 0 {
		for i := start; i < end; i++ {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), prefix) {
				lines[i] = newLine
				return lines
			}
		}
		lines = slices.Insert(lines, end, newLine)
		return lines
	}

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
