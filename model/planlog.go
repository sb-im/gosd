package model

type PlanLog struct {
	ID          int64             `json:"id"`
	LogID       int64             `json:"log_id"`
	PlanID      int64             `json:"plan_id"`
	Attachments map[string]string `json:"attachments"`
	Extra       map[string]string `json:"extra"`
	RecordTime
}

func NewPlanLog() *PlanLog {
	return &PlanLog{
		Attachments: make(map[string]string),
		Extra:       make(map[string]string),
	}
}

type PlanLogs []*PlanLog
