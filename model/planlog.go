package model

type PlanLog struct {
	ID     int64             `json:"id"`
	LogID  int64             `json:"job_id"`
	PlanID int64             `json:"plan_id"`
	Files  map[string]string `json:"files"`
	Extra  map[string]string `json:"extra"`
	RecordTime
}

func NewPlanLog() *PlanLog {
	return &PlanLog{
		Files: make(map[string]string),
		Extra: make(map[string]string),
	}
}

type PlanLogs []*PlanLog
