package model

type JudgeRequest struct {
	SubmitID    string   `json:"submitID"`
	TaskID      string   `json:"taskID"`
	LanguageID  string   `json:"languageID"`
	TestcaseIDs []string `json:"testcaseIDs"`
}

type JudgeResponse struct {
	Status          Status           `json:"status"`
	CompileMessage  *string          `json:"compileMessage"` // nil if Status is not CE
	ErrorMessage    *string          `json:"errorMessage"`   // nil if Status is not IE
	TestcaseResults []TestcaseResult `json:"testcaseResults"`
}

type TestcaseResult struct {
	ID              string `json:"id"`
	Status          Status `json:"status"`
	ExecutionTime   int64  `json:"executionTime"`   // in milliseconds
	ExecutionMemory int64  `json:"executionMemory"` // in kilobytes
}

type Status string

const (
	StatusAC  Status = "AC"
	StatusWA  Status = "WA"
	StatusRE  Status = "RE"
	StatusTLE Status = "TLE"
	StatusMLE Status = "MLE"
	StatusOLE Status = "OLE"
	StatusCE  Status = "CE"
	StatusIE  Status = "IE"
)