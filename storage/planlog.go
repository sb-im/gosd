package storage

import (
	"database/sql"
	"fmt"

	"sb.im/gosd/model"

	"github.com/lib/pq/hstore"
)

func (s *Storage) CreatePlanLog(item *model.PlanLog) (err error) {
	attachments := hstore.Hstore{Map: make(map[string]sql.NullString)}

	if len(item.Attachments) > 0 {
		for key, value := range item.Attachments {
			attachments.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	extra := hstore.Hstore{Map: make(map[string]sql.NullString)}

	if len(item.Extra) > 0 {
		for key, value := range item.Extra {
			extra.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	query := `
		INSERT INTO plan_logs
			(log_id, plan_id, attachments, extra, create_at, update_at)
		VALUES
			(
				(
					SELECT
						CASE WHEN MAX(log_id) IS NULL THEN 1 ELSE MAX(log_id) + 1 END
					FROM
						plan_logs WHERE plan_id=$1
				), $1, $2, $3, now(), now()
			)
		RETURNING
			id, log_id, plan_id
	`

	err = s.db.QueryRow(query, item.PlanID, attachments, extra).Scan(
		&item.ID,
		&item.LogID,
		&item.PlanID,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to create plan_log: %v`, err)
	}

	return nil
}
