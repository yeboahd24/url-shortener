package sqlc

import (
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
)

// UUIDToNullable converts a *uuid.UUID to pgtype.UUID
func UUIDToNullable(id *uuid.UUID) pgtype.UUID {
    if id == nil {
        return pgtype.UUID{Valid: false}
    }
    return pgtype.UUID{Bytes: *id, Valid: true}
}