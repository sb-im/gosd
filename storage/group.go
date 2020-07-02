package storage

import (
	"database/sql"
	"fmt"

	"sb.im/gosd/model"

	"github.com/lib/pq/hstore"
)

func (s *Storage) CreateGroup(user *model.Group) (err error) {
	extra := hstore.Hstore{Map: make(map[string]sql.NullString)}

	if len(user.Extra) > 0 {
		for key, value := range user.Extra {
			extra.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	query := `
		INSERT INTO groups
			(name, extra)
		VALUES
			($1, $2)
		RETURNING
			id, name
	`

	err = s.db.QueryRow(query, user.Name, extra).Scan(
		&user.ID,
		&user.Name,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to create group: %v`, err)
	}

	return nil
}
