package config

type MainConfig struct {
	General GeneralConfig `yaml:"general"`
	Git     GitConfig     `yaml:"git"`
	Include []string      `yaml:"include"`
}

type GeneralConfig struct {
	DefaultProfile string `yaml:"default_profile"`
	AutoSwitch     bool   `yaml:"auto_switch"`
	SSHConfigPath  string `yaml:"ssh_config_path"`
}

type GitConfig struct {
	AliasPrefix string `yaml:"alias_prefix"`
}

type Profile struct {
	Profile ProfileMeta `yaml:"profile"`
	User    UserConfig  `yaml:"user"`
	SSH     *SSHConfig  `yaml:"ssh,omitempty"`
	GPG     *GPGConfig  `yaml:"gpg,omitempty"`
	Rules   RulesConfig `yaml:"rules"`
	Include []string    `yaml:"include,omitempty"`
}

type ProfileMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
}

type UserConfig struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type SSHConfig struct {
	KeyFile   string   `yaml:"key_file"`
	Hosts     []string `yaml:"hosts"`
	HostAlias string   `yaml:"host_alias,omitempty"`
}

type GPGConfig struct {
	SigningKey  string `yaml:"signing_key,omitempty"`
	SignCommits bool   `yaml:"sign_commits"`
}

type RulesConfig struct {
	Directory []string `yaml:"directory"`
	URL       []string `yaml:"url"`
}

func DefaultMainConfig() MainConfig {
	return MainConfig{
		General: GeneralConfig{
			DefaultProfile: "personal",
			AutoSwitch:     true,
			SSHConfigPath:  "~/.ssh/config",
		},
		Git: GitConfig{
			AliasPrefix: "sw",
		},
		Include: []string{"profiles/*.yaml"},
	}
}
