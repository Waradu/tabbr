package shell

import _ "embed"

//go:embed tabbr.zsh
var Zsh string

//go:embed tabbr.bash
var Bash string

//go:embed tabbr.ps1
var PowerShell string
