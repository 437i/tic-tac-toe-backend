package datasource

import "errors"

var ErrBuildConnString = errors.New("missing required environment variables for DB connection")
