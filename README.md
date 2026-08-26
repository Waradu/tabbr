# Tabbr

Fuzzy autocomplete for commands you actually use.

```text
brd<Tab>  →  bun run dev
```

Tabbr learns commands that finish successfully. Tab expands the best match, repeated Tab cycles through matches, and
normal shell completion remains available.

## Things to note

- Tabbr only stores successfully completed commands that contain at least six characters.
- Completion requires a whitespace-free query of at least two characters.
- Commands not used within the last month are ignored.

> [!CAUTION]
> The built-in exclusions are not exhaustive. Commands are stored unencrypted in a local SQLite database, so add
exclusion patterns for anything you do not want stored.

## Install

```sh
go install github.com/waradu/tabbr/cmd/tabbr@latest
```

Ensure `$(go env GOPATH)/bin` is in your `PATH`.

## Shell setup

### Zsh

Add to `~/.zshrc`:

```zsh
eval "$(tabbr init zsh)"
```

### Bash

Bash 4.4 or newer is required.

Add to `~/.bashrc`:

```bash
eval "$(tabbr init bash)"
```

### PowerShell

PowerShell integration requires PowerShell 5.1 or newer and PSReadLine 2.0 or newer (should be installed by default).

Add to your PowerShell profile:

```powershell
Invoke-Expression (& { (tabbr init pwsh | Out-String) })
```

## Commands

### Shell integration

#### `tabbr init <shell>`

Print the adapter for `zsh`, `bash`, or `pwsh`.

### Command history

#### `tabbr add "<command>"`

Add a command manually.

#### `tabbr remove "<command>"`

Remove a stored command.

#### `tabbr query "<query>"`

Print matching commands in ranked order.

#### `tabbr list`

List stored commands with their score and last-used time.

#### `tabbr path`

Print the absolute path to the local SQLite database.

### Exclusions

Exclusions use case-insensitive glob patterns such as `*--token=*`, not regular expressions.
Built-in exclusions cover `cd *`, `ls *` and common sensitive values such as `API_KEY` and `--api-key`.

#### `tabbr exclude add "<pattern>"`

Add an exclusion and remove existing matching commands.

#### `tabbr exclude remove "<pattern>"`

Remove an exclusion.

#### `tabbr exclude list`

List all exclusion patterns.
