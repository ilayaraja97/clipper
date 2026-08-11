package ai

type EngineExecOutput struct {
	Command     string `json:"cmd"`
	Explanation string `json:"exp"`
	Executable  bool   `json:"exec"`
}

func (eo EngineExecOutput) GetCommand() string {
	return eo.Command
}

func (eo EngineExecOutput) GetExplanation() string {
	return eo.Explanation
}

func (eo EngineExecOutput) IsExecutable() bool {
	return eo.Executable
}
