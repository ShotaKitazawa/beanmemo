//go:build tools

package hack

import (
	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
	_ "github.com/ogen-go/ogen/cmd/ogen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
