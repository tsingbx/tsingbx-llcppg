package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgParse dbgFlags = 1 << iota
	DbgVisitTop
	DbgProcess
	DbgGetCurFile
	DbgFlagAll = DbgParse
)

func SetDebugParse() {
	flags |= DbgParse
}

func GetDebugParse() bool {
	return flags&DbgParse != 0
}

func SetDebugAll() {
	flags = DbgFlagAll
}

func SetDebugVisitTop() {
	flags |= DbgVisitTop
}

func GetDebugVisitTop() bool {
	return flags&DbgVisitTop != 0
}

func SetDebugProcess() {
	flags |= DbgProcess
}

func GetDebugProcess() bool {
	return flags&DbgProcess != 0
}

func SetDebugGetCurFile() {
	flags |= DbgGetCurFile
}

func GetDebugGetCurFile() bool {
	return flags&DbgGetCurFile != 0
}
