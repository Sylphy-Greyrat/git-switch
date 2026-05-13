module github.com/sylphy/git-switch/cli

go 1.22

require (
	github.com/sylphy/git-switch/core v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace github.com/sylphy/git-switch/core => ../core
