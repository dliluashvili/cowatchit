package constants

import "time"

type contextKey string

const SessionContextKey contextKey = "session"

type validatedKey string

const ValidatedContextKey validatedKey = "validated"

var SessionDuration = 24 * time.Hour
