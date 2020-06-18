package storage

import (
	"database/sql"
	"fmt"

	"sb.im/gosd/model"

	"github.com/lib/pq/hstore"
)

func (s *Storage) CreatePlan(plan *model.Plan) (err error) {
	attachments := hstore.Hstore{Map: make(map[string]sql.NullString)}

	if len(plan.Attachments) > 0 {
		for key, value := range plan.Attachments {
			attachments.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	extra := hstore.Hstore{Map: make(map[string]sql.NullString)}

	if len(plan.Extra) > 0 {
		for key, value := range plan.Extra {
			extra.Map[key] = sql.NullString{String: value, Valid: true}
		}
	}

	query := `
	INSERT INTO plans
	(name, description, node_id, attachments, extra, create_at, update_at)
	VALUES
	($1, $2, $3, $4, $5, now(), now())
	RETURNING
	id, name, description, node_id
	`

	err = s.db.QueryRow(query, plan.Name, plan.Description, plan.NodeID, attachments, extra).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.NodeID,
	)

	if err != nil {
		return fmt.Errorf(`store: unable to create plan: %v`, err)
	}

	return nil
}

// Plans returns all plans.
func (s *Storage) Plans() (model.Plans, error) {
	query := `
		SELECT
			id,
			name,
			description,
			node_id,
			attachments,
			extra
		FROM
			plans
		ORDER BY id ASC
	`

	rows, err := s.db.Query(query)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch plans: %v`, err)
	}

	var plans model.Plans

	for rows.Next() {
		var attachments hstore.Hstore
		var extra hstore.Hstore
		plan := model.NewPlan()
		err := rows.Scan(
			&plan.ID,
			&plan.Name,
			&plan.Description,
			&plan.NodeID,
			&attachments,
			&extra,
		)

		if err != nil {
			return nil, fmt.Errorf(`store: unable to fetch plans row: %v`, err)
		}

		for key, value := range attachments.Map {
			if value.Valid {
				plan.Attachments[key] = value.String
			}
		}

		for key, value := range extra.Map {
			if value.Valid {
				plan.Extra[key] = value.String
			}
		}

		plans = append(plans, plan)
	}

	return plans, nil
}

// PlanByID finds a plan by the ID.
func (s *Storage) PlanByID(planID int64) (*model.Plan, error) {
	query := `
		SELECT
			id,
			name,
			description,
			node_id,
			attachments,
			extra
		FROM
			plans
		WHERE
			id = $1
	`

	return s.fetchPlan(query, planID)
}

func (s *Storage) fetchPlan(query string, args ...interface{}) (*model.Plan, error) {
	var attachments hstore.Hstore
	var extra hstore.Hstore
	plan := model.NewPlan()

	err := s.db.QueryRow(query, args...).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Description,
		&plan.NodeID,
		&attachments,
		&extra,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf(`store: unable to fetch user: %v`, err)
	}

	for key, value := range attachments.Map {
		if value.Valid {
			plan.Attachments[key] = value.String
		}
	}

	for key, value := range extra.Map {
		if value.Valid {
			plan.Extra[key] = value.String
		}
	}

	return plan, nil
}
