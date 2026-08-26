#!/usr/bin/env bash

[[ $- == *i* ]] || return

if (( BASH_VERSINFO[0] < 4 || BASH_VERSINFO[0] == 4 && BASH_VERSINFO[1] < 4 )); then
	printf 'tabbr: Bash 4.4 or newer is required\n' >&2
	return
fi

_tabbr_complete() {
	local match
	COMPREPLY=()

	(( COMP_POINT == ${#COMP_LINE} )) || return 0
	[[ $COMP_LINE != *[[:space:]]* ]] || return 0

	while IFS= read -r match; do
		[[ -n $match ]] && COMPREPLY+=("$match")
	done < <(command tabbr query "$COMP_LINE" 2>/dev/null)

	if (( ${#COMPREPLY[@]} > 0 )); then
		compopt -o noquote -o nosort -o nospace
	fi
}

_tabbr_record_prompt() {
	local exit_status=$?
	local entry history_id executed
	local HISTTIMEFORMAT=

	entry="$(builtin history 1)"
	if [[ ! $entry =~ ^[[:space:]]*([0-9]+)[[:space:]]+(.*)$ ]]; then
		return 0
	fi

	history_id=${BASH_REMATCH[1]}
	executed=${BASH_REMATCH[2]}

	if [[ -z ${_tabbr_history_ready:-} ]]; then
		_tabbr_history_ready=1
		_tabbr_last_history_id=$history_id
		return 0
	fi

	[[ $history_id != "${_tabbr_last_history_id:-}" ]] || return 0
	_tabbr_last_history_id=$history_id

	(( exit_status == 0 )) || return 0

	command tabbr add "$executed" >/dev/null 2>&1
	return 0
}

complete -I -F _tabbr_complete -o bashdefault -o default
bind '"\C-i": menu-complete'
bind '"\e[Z": menu-complete-backward'

if [[ ${PROMPT_COMMAND:-} != *'_tabbr_record_prompt'* ]]; then
	PROMPT_COMMAND="_tabbr_record_prompt${PROMPT_COMMAND:+; $PROMPT_COMMAND}"
fi
