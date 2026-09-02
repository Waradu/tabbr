#!/usr/bin/env zsh
# shellcheck disable=SC2296
# shellcheck disable=SC2034

typeset -ga _tabbr_matches
typeset -gi _tabbr_index=0
typeset -g _tabbr_fallback_widget
typeset -g _tabbr_pending_command

_tabbr_reset() {
	_tabbr_matches=()
	_tabbr_index=0
}

_tabbr_fallback() {
	local widget="${1:-$_tabbr_fallback_widget}"
	_tabbr_reset
	zle "$widget"
}

_tabbr_cycle() {
	local direction="$1"
	local fallback="$2"
	local output

	if (( ${#_tabbr_matches} > 0 )) &&
		[[ "$BUFFER" == "${_tabbr_matches[$((_tabbr_index + 1))]}" ]] &&
		(( CURSOR == ${#BUFFER} )); then
		_tabbr_index=$(( (_tabbr_index + direction + ${#_tabbr_matches}) % ${#_tabbr_matches} ))
		BUFFER="${_tabbr_matches[$((_tabbr_index + 1))]}"
		CURSOR=${#BUFFER}
		POSTDISPLAY=
		return
	fi

	_tabbr_reset

	if (( CURSOR != ${#BUFFER} )); then
		_tabbr_fallback "$fallback"
		return
	fi

	output="$(command tabbr query "$BUFFER" 2>/dev/null)" || {
		_tabbr_fallback "$fallback"
		return
	}
	if [[ -z "$output" ]]; then
		_tabbr_fallback "$fallback"
		return
	fi
	_tabbr_matches=("${(@f)output}")
	(( direction < 0 )) && _tabbr_index=$((${#_tabbr_matches} - 1))

	BUFFER="${_tabbr_matches[$((_tabbr_index + 1))]}"
	CURSOR=${#BUFFER}
	POSTDISPLAY=
}

_tabbr_complete() {
	_tabbr_cycle 1 "$_tabbr_fallback_widget"
}

_tabbr_complete_reverse() {
	_tabbr_cycle -1 reverse-menu-complete
}

_tabbr_record_preexec() {
	_tabbr_pending_command="$1"
}

_tabbr_record_precmd() {
	local exit_status=$?
	local executed="$_tabbr_pending_command"
	_tabbr_pending_command=

	(( exit_status == 0 )) || return 0

	command tabbr add "$executed" >/dev/null 2>&1
	return 0
}

if [[ -z "$_tabbr_fallback_widget" ]]; then
	typeset -a _tabbr_binding
	_tabbr_binding=("${(z)$(bindkey '^I')}")
	_tabbr_fallback_widget="${_tabbr_binding[2]:-expand-or-complete}"
	unset _tabbr_binding
fi

zle -N _tabbr_complete
zle -N _tabbr_complete_reverse
bindkey '^I' _tabbr_complete
bindkey '^[[Z' _tabbr_complete_reverse

autoload -Uz add-zsh-hook
add-zsh-hook -d preexec _tabbr_record_preexec 2>/dev/null
add-zsh-hook -d precmd _tabbr_record_precmd 2>/dev/null
add-zsh-hook preexec _tabbr_record_preexec
add-zsh-hook precmd _tabbr_record_precmd
