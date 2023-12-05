package internal

type ACTION string

const (
	GETSOURCECODE      ACTION = "getSourceCode"
	LISTEXPLORER       ACTION = "listExplorer"
	GETSOURCECODEBYURL ACTION = "getSourceCodeByUrl"
	SHOWVERSION        ACTION = "showVersion"
	NONE               ACTION = "none"
)
