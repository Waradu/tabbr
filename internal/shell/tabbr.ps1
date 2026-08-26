if (-not (Get-Module PSReadLine)) {
	Import-Module PSReadLine -ErrorAction SilentlyContinue
}

if (-not (Get-Module PSReadLine)) {
	Write-Warning 'tabbr: PSReadLine is required'
	return
}

function global:_tabbr_reset {
	$global:_tabbr_matches = @()
	$global:_tabbr_index = 0
}

function global:_tabbr_fallback {
	param(
		[System.ConsoleKeyInfo] $Key,
		[object] $Argument,
		[int] $Direction
	)

	_tabbr_reset
	if ($Direction -lt 0) {
		[Microsoft.PowerShell.PSConsoleReadLine]::TabCompletePrevious($Key, $Argument)
	} else {
		[Microsoft.PowerShell.PSConsoleReadLine]::TabCompleteNext($Key, $Argument)
	}
}

function global:_tabbr_cycle {
	param(
		[System.ConsoleKeyInfo] $Key,
		[object] $Argument,
		[int] $Direction
	)

	$line = $null
	$cursor = 0
	[Microsoft.PowerShell.PSConsoleReadLine]::GetBufferState([ref] $line, [ref] $cursor)

	if ($global:_tabbr_matches.Count -gt 0 -and
		$line -ceq $global:_tabbr_matches[$global:_tabbr_index] -and
		$cursor -eq $line.Length) {
		$global:_tabbr_index = (
			$global:_tabbr_index + $Direction + $global:_tabbr_matches.Count
		) % $global:_tabbr_matches.Count
		$match = $global:_tabbr_matches[$global:_tabbr_index]
		[Microsoft.PowerShell.PSConsoleReadLine]::Replace(0, $line.Length, $match)
		[Microsoft.PowerShell.PSConsoleReadLine]::SetCursorPosition($match.Length)
		return
	}

	_tabbr_reset

	if ($cursor -ne $line.Length -or $line -match '\s') {
		_tabbr_fallback $Key $Argument $Direction
		return
	}

	$tabbrMatches = @(& tabbr query $line 2>$null)
	$querySucceeded = $?
	$tabbrMatches = @($tabbrMatches | Where-Object { -not [string]::IsNullOrEmpty($_) })
	if (-not $querySucceeded -or $tabbrMatches.Count -eq 0) {
		_tabbr_fallback $Key $Argument $Direction
		return
	}

	$global:_tabbr_matches = $tabbrMatches
	if ($Direction -lt 0) {
		$global:_tabbr_index = $tabbrMatches.Count - 1
	}

	$match = $global:_tabbr_matches[$global:_tabbr_index]
	[Microsoft.PowerShell.PSConsoleReadLine]::Replace(0, $line.Length, $match)
	[Microsoft.PowerShell.PSConsoleReadLine]::SetCursorPosition($match.Length)
}

function global:_tabbr_complete {
	param($Key, $Argument)
	_tabbr_cycle $Key $Argument 1
}

function global:_tabbr_complete_reverse {
	param($Key, $Argument)
	_tabbr_cycle $Key $Argument -1
}

function global:_tabbr_record_prompt {
	param([bool] $Succeeded)

	$entry = Get-History -Count 1 -ErrorAction SilentlyContinue
	if ($null -eq $entry -or $entry.Id -eq $global:_tabbr_last_history_id) {
		return
	}

	$global:_tabbr_last_history_id = $entry.Id
	if (-not $Succeeded) {
		return
	}

	$executed = $entry.CommandLine
	$lastExitCode = $global:LASTEXITCODE
	& tabbr add $executed *> $null
	$global:LASTEXITCODE = $lastExitCode
}

_tabbr_reset
$lastHistory = Get-History -Count 1 -ErrorAction SilentlyContinue
$global:_tabbr_last_history_id = if ($null -eq $lastHistory) { -1 } else { $lastHistory.Id }

Set-PSReadLineKeyHandler -Chord Tab -ScriptBlock ${function:_tabbr_complete} `
	-BriefDescription 'TabbrComplete' `
	-Description 'Complete using commands learned by tabbr'
Set-PSReadLineKeyHandler -Chord Shift+Tab -ScriptBlock ${function:_tabbr_complete_reverse} `
	-BriefDescription 'TabbrCompleteReverse' `
	-Description 'Complete backwards using commands learned by tabbr'

if (-not $global:_tabbr_prompt_installed) {
	$global:_tabbr_original_prompt = ${function:prompt}
	function global:prompt {
		$succeeded = $?
		$lastExitCode = $global:LASTEXITCODE
		_tabbr_record_prompt $succeeded
		$global:LASTEXITCODE = $lastExitCode
		& $global:_tabbr_original_prompt
	}
	$global:_tabbr_prompt_installed = $true
}
