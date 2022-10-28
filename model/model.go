package model

type JudgeRequest struct {
	SubmitID   int        `json:"submit_id"`
	TaskID     int        `json:"task_id"`
	LanguageID string     `json:"language_id"`
	Testcases  []Testcase `json:"testcases"`
}

type Testcase struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JudgeResponse struct {
	Status          Status           `json:"status"`
	ExecutionTime   int64            `json:"execution_time"`   // max usage
	ExecutionMemory int64            `json:"execution_memory"` // max usage
	CompileMessage  *string          `json:"compile_message"`  // nil if Status is not CE
	ErrorMessage    *string          `json:"error_message"`    // nil if Status is not IE
	TestcaseResults []TestcaseResult `json:"testcase_results"`
}

type TestcaseResult struct {
	ID              int    `json:"id"`
	Status          Status `json:"status"`
	ExecutionTime   int64  `json:"execution_time"`   // in milliseconds
	ExecutionMemory int64  `json:"execution_memory"` // in kilobytes
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
