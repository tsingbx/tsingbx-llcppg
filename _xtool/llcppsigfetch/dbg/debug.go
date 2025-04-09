package dbg

type dbgFlags = int

var flags dbgFlags

const (
	DbgParse dbgFlags = 1 << iota
	DbgVisitTop
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
