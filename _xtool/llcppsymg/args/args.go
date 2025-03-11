package args

import (
	"strconv"
	"strings"
)

const LLCPPG_CFG = "llcppg.cfg"
const LLCPPG_SYMB = "llcppg.symb.json"
const LLCPPG_SIGFETCH = "llcppg.sigfetch.json"
const LLCPPG_PUB = "llcppg.pub"

type Args struct {
	Help                 bool
	Verbose              bool
	VerboseSigfetchParse bool //-vsp llcppsigfetch parse.go
	VerboseParseIsMethod bool //-vpim
	UseStdin             bool
	ParseAllComment      bool
	CfgFile              string
}

func ParseArgs(args []string, defaultCfgFile string, swflags map[string]bool) (*Args, []string) {
	result := &Args{}
	filteredArgs := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			switch arg {
			case "-h", "--help":
				result.Help = true
				continue
			case "-v":
				result.Verbose = true
				continue
			case "-parse-all-comment":
				result.ParseAllComment = true
				continue
			case "-vpim":
				result.VerboseParseIsMethod = true
				continue
			case "-vsp":
				result.VerboseSigfetchParse = true
				continue
			case "-":
				result.UseStdin = true
				continue
			default:
				if hasArg, ok := swflags[arg]; ok {
					if hasArg {
						filteredArgs = append(filteredArgs, arg)
						for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
							filteredArgs = append(filteredArgs, args[i+1])
							i++
						}
						continue
					}
				}
				filteredArgs = append(filteredArgs, arg)
			}
		} else if result.CfgFile == "" {
			result.CfgFile = arg
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	if result.CfgFile == "" {
		result.CfgFile = defaultCfgFile
	}
	return result, filteredArgs
}

func BoolArg(arg string, defaultValue bool) bool {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return defaultValue
	}
	value, err := strconv.ParseBool(parts[1])
	if err != nil {
		return defaultValue
	}
	return value
}

func StringArg(arg string, defaultValue string) string {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return defaultValue
	}
	return parts[1]
}
