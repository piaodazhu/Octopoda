package protocols

type Result struct {
	// Result code. Reserved
	Rcode int

	// Result msg string. OK, or Error Reason
	Rmsg string

	// For script or command execution. Script output.
	Output string

	// For version control. Version hash code.
	Version string

	// For version control. Modified flag.
	Modified bool
}

type ExecutionResult struct {
	Name                  string
	Code                  ExecStatus
	CommunicationErrorMsg string
	ProcessErrorMsg       string
	Result                string
	ResultCompatible      string
}

type ExecStatus int

const (
	ExecOK ExecStatus = iota
	ExecCommunicationError
	ExecProcessError
)

type ExecutionResults []ExecutionResult
type ExecutionResultsText struct {
	Results            ExecutionResults
	Total              int
	Success            int
	CommunicationError int
	ProcessError       int
}

func (results *ExecutionResults) ToText() *ExecutionResultsText {
	res := &ExecutionResultsText{
		Results: *results,
		Total: len(*results),
		Success: 0,
		CommunicationError: 0,
		ProcessError: 0,
	}
	for _, result := range *results {
		if result.Code == 0 {
			res.Success++
		} else if result.Code == ExecCommunicationError {
			res.CommunicationError++
		} else if result.Code == ExecProcessError {
			res.ProcessError++
		}
	}
	return res
}
