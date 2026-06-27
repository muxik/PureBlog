// Package migrations embeds the goose SQL migrations so they ship inside the
// binary and run automatically on startup.
package migrations

import "embed"

// FS holds the embedded *.sql migration files.
//
//go:embed *.sql
var FS embed.FS
