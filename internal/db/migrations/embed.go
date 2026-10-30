// Package migrations embeds the SQL migration files into the compiled
// binary, so the server can apply them on startup without needing the
// .sql files to exist next to it (important later for Docker/Railway).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
