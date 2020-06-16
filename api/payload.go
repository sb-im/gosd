package api

import (
	"encoding/json"
	"fmt"
	"io"

	"sb.im/gosd/model"
)

func decodePlanCreationPayload(r io.ReadCloser) (*model.Plan, error) {
	defer r.Close()

	var plan model.Plan
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&plan); err != nil {
		return nil, fmt.Errorf("Unable to decode plan modification JSON object: %v", err)
	}
	return &plan, nil
}
