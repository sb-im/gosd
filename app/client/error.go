package client

type errMsg struct {
	status string
	Errors string `json:"error"`
}

func (e *errMsg) Error() string {
	return e.status + " " + e.Errors
}
