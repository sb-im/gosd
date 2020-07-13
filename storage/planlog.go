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

// PlanLogs returns all logs.
func (s *Storage) PlanLogs(planID int64) (model.PlanLogs, error) {
	query := `
		SELECT
			id,
			log_id,
			attachments,
			extra
		FROM
			plan_logs
		WHERE
			plan_id = $1
		ORDER BY log_id ASC
	`

	return s.fetchPlanLogs(query, planID)
}

func (s *Storage) fetchPlanLogs(query string, args ...interface{}) (model.PlanLogs, error) {
	rows, err := s.db.Query(query, args...)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch plan_logs: %v`, err)
	}

	var logs model.PlanLogs

	for rows.Next() {
		var attachments hstore.Hstore
		var extra hstore.Hstore
		log := model.NewPlanLog()
		err := rows.Scan(
			&log.ID,
			&log.LogID,
			&attachments,
			&extra,
		)

		if err != nil {
			return nil, fmt.Errorf(`store: unable to fetch plan_logs row: %v`, err)
		}

		for key, value := range attachments.Map {
			if value.Valid {
				log.Attachments[key] = value.String
			}
		}

		for key, value := range extra.Map {
			if value.Valid {
				log.Extra[key] = value.String
			}
		}

		logs = append(logs, log)
	}

	return logs, nil
}
