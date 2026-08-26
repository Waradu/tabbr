package store

import (
	"database/sql"
	"strings"
	"time"
	"unicode/utf8"
)

func Add(db *sql.DB, command string) error {
	if strings.ContainsAny(command, "\r\n") {
		return nil
	}

	command = strings.TrimSpace(command)

	fields := strings.Fields(command)
	if len(fields) > 0 && fields[0] == "tabbr" {
		return nil
	}

	if utf8.RuneCountInString(command) <= 5 {
		return nil
	}

	var excluded bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM exclusions
			WHERE lower(?) GLOB lower(pattern)
		)
	`, command).Scan(&excluded)
	if err != nil || excluded {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO commands (command, use_count, last_used_at)
		VALUES (?, 1, ?)
		ON CONFLICT(command) DO UPDATE SET
			use_count = use_count + 1,
			last_used_at = excluded.last_used_at
	`, command, time.Now().Unix())

	return err
}

func Remove(db *sql.DB, command string) error {
	_, err := db.Exec(`
		DELETE FROM commands
		WHERE command = ?
	`, command)

	return err
}

func List(db *sql.DB) ([]Command, error) {
	rows, err := db.Query(`
		SELECT command, use_count, last_used_at FROM commands
		WHERE last_used_at > unixepoch('now', '-1 month')
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []Command

	for rows.Next() {
		var command Command

		if err := rows.Scan(
			&command.Text,
			&command.UseCount,
			&command.LastUsed,
		); err != nil {
			return nil, err
		}

		commands = append(commands, command)
	}

	return commands, rows.Err()
}

func AddExclusion(db *sql.DB, pattern string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`
		INSERT INTO exclusions (pattern, created_at)
		VALUES (?, ?)
		ON CONFLICT(pattern) DO NOTHING
	`, pattern, time.Now().Unix())
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM commands
		WHERE lower(command) GLOB lower(?)
	`, pattern)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func RemoveExclusion(db *sql.DB, pattern string) error {
	_, err := db.Exec(`
		DELETE FROM exclusions
		WHERE pattern = ?
	`, pattern)

	return err
}

func ListExclusions(db *sql.DB) ([]Exclusion, error) {
	rows, err := db.Query(`
			SELECT pattern, created_at FROM exclusions
			ORDER BY created_at DESC, pattern
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exclusions []Exclusion

	for rows.Next() {
		var exclusion Exclusion

		if err := rows.Scan(
			&exclusion.Pattern,
			&exclusion.CreatedAt,
		); err != nil {
			return nil, err
		}

		exclusions = append(exclusions, exclusion)
	}

	return exclusions, rows.Err()
}
