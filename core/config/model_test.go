package config

import "testing"

func TestProfileModelFields(t *testing.T) {
	profile := Profile{
		Profile: ProfileMeta{
			Name:        "personal",
			Description: "Personal GitHub account",
		},
		User: UserConfig{
			Name:  "Sylphy",
			Email: "sylphy@example.com",
		},
		SSH: &SSHConfig{
			KeyFile:   "~/.ssh/id_rsa",
			Hosts:     []string{"github.com"},
			HostAlias: "github.com-personal",
		},
		Rules: RulesConfig{
			Directory: []string{"~/projects/personal/*"},
			URL:       []string{"github.com:sylphy/*"},
		},
	}

	if profile.Profile.Name != "personal" {
		t.Fatalf("expected profile name personal, got %q", profile.Profile.Name)
	}
	if profile.User.Email != "sylphy@example.com" {
		t.Fatalf("expected email sylphy@example.com, got %q", profile.User.Email)
	}
	if profile.SSH == nil || profile.SSH.HostAlias != "github.com-personal" {
		t.Fatalf("expected SSH host alias github.com-personal, got %#v", profile.SSH)
	}
}
