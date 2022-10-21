package model

type JudgeRequest struct {
	SubmitID    uint64   `json:"submitID"`
	LanguageID  string   `json:"languageID"`
	TestcaseIDs []uint64 `json:"testcaseIDs"`
}

type JudgeResponse struct {
	Status          Status           `json:"status"`
	CompileMessage  *string          `json:"compileMessage"` // nil if Status is not CE
	ErrorMessage    *string          `json:"errorMessage"`   // nil if Status is not IE
	TestcaseResults []TestcaseResult `json:"testcaseResults"`
}

type TestcaseResult struct {
	ID              uint64 `json:"id"`
	Status          Status `json:"status"`
	ExecutionTime   int64  `json:"executionTime"`   // in milliseconds
	ExecutionMemory int64  `json:"executionMemory"` // in killobytes
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
