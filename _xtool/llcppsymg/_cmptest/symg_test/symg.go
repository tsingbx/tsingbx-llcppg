package main

import (
	"fmt"
	"os"

	"github.com/goplus/llcppg/_xtool/llcppsymg/parse"
	"github.com/goplus/llcppg/_xtool/llcppsymg/symbol"
	"github.com/goplus/llgo/xtool/nm"
)

func main() {
	TestParseHeaderFile()
}
func TestParseHeaderFile() {
	testCases := []struct {
		name            string
		content         string
		isCpp           bool
		prefixes        []string
		dylibSymbols    []*nm.Symbol
		symbFileContent string
	}{
		{
			name: "inireader",
			content: `
#define INI_API __attribute__((visibility("default")))
class INIReader {
  public:
    __attribute__((visibility("default"))) explicit INIReader(const char *filename);
    INI_API explicit INIReader(const char *buffer, long buffer_size);
    ~INIReader();
    INI_API int ParseError() const;
	INI_API const char * Get(const char *section, const char *name,
						const char *default_value) const;
  private:
    static const char * MakeKey(const char *section, const char *name);
};
            `,
			isCpp:    true,
			prefixes: []string{"INI"},
			dylibSymbols: []*nm.Symbol{
				{Name: symbol.AddSymbolPrefixUnder("ZN9INIReaderC1EPKc", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZN9INIReaderC1EPKcl", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZN9INIReaderD1Ev", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZNK9INIReader10ParseErrorEv", true)},
				{Name: symbol.AddSymbolPrefixUnder("ZNK9INIReader3GetEPKcS1_S1_", true)},
			},
			symbFileContent: `
[{
		"mangle":       "_ZN9INIReaderC1EPKc",
		"c++":  "INIReader::INIReader(const char *)",
		"go":   "(*Reader).Init"
}, {
		"mangle":       "_ZN9INIReaderC1EPKcl",
		"c++":  "INIReader::INIReader(const char *, long)",
		"go":   "(*Reader).Init__1"
}, {
		"mangle":       "_ZN9INIReaderD1Ev",
		"c++":  "INIReader::~INIReader()",
		"go":   "(*Reader).Dispose"
}, {
		"mangle":       "_ZNK9INIReader10ParseErrorEv",
		"c++":  "INIReader::ParseError()",
		"go":   "(*Reader).ModifyedParseError"
}]`,
		},
		{
			name: "lua",
			content: `
typedef struct lua_State lua_State;

LUA_API int(lua_error)(lua_State *L);

LUA_API int(lua_next)(lua_State *L, int idx);

LUA_API void(lua_concat)(lua_State *L, int n);
LUA_API void(lua_len)(lua_State *L, int idx);

LUA_API long unsigned int(lua_stringtonumber)(lua_State *L, const char *s);

LUA_API void(lua_setallocf)(lua_State *L, lua_Alloc f, void *ud);

LUA_API void(lua_toclose)(lua_State *L, int idx);
LUA_API void(lua_closeslot)(lua_State *L, int idx);
            `,
			isCpp:    false,
			prefixes: []string{"lua_"},
			dylibSymbols: []*nm.Symbol{
				{Name: symbol.AddSymbolPrefixUnder("lua_error", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_next", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_concat", false)},
				{Name: symbol.AddSymbolPrefixUnder("lua_stringtonumber", false)},
			},
		},
		{
			name: "cjson",
			content: `
		#define CJSON_PUBLIC(type) type
		#include <stddef.h>
		/* The cJSON structure: */
		typedef struct cJSON {
		    /* next/prev allow you to walk array/object chains. Alternatively, use GetArraySize/GetArrayItem/GetObjectItem */
		    struct cJSON *next;
		    struct cJSON *prev;
		    /* An array or object item will have a child pointer pointing to a chain of the items in the array/object. */
		    struct cJSON *child;

		    /* The type of the item, as above. */
		    int type;

		    /* The item's string, if type==cJSON_String  and type == cJSON_Raw */
		    char *valuestring;
		    /* writing to valueint is DEPRECATED, use cJSON_SetNumberValue instead */
		    int valueint;
		    /* The item's number, if type==cJSON_Number */
		    double valuedouble;

		    /* The item's name string, if this item is the child of, or is in the list of subitems of an object. */
		    char *string;
		} cJSON;
		/* Render a cJSON entity to text for transfer/storage. */
		CJSON_PUBLIC(char *) cJSON_Print(const cJSON *item);
		CJSON_PUBLIC(cJSON *) cJSON_ParseWithLength(const char *value, size_t buffer_length);
		/* Delete a cJSON entity and all subentities. */
		CJSON_PUBLIC(void) cJSON_Delete(cJSON *item);
					`,
			isCpp:    false,
			prefixes: []string{"cJSON_"},
			dylibSymbols: []*nm.Symbol{
				{Name: symbol.AddSymbolPrefixUnder("cJSON_Print", false)},
				{Name: symbol.AddSymbolPrefixUnder("cJSON_ParseWithLength", false)},
				{Name: symbol.AddSymbolPrefixUnder("cJSON_Delete", false)},
			},
		},
	}

	for _, tc := range testCases {
		fmt.Printf("=== Test Case: %s ===\n", tc.name)
		headerSymbolMap, err := parse.ParseHeaderFile([]string{tc.content}, tc.prefixes, tc.isCpp, true)
		if err != nil {
			fmt.Println("Error:", err)
		}
		tmpFile, err := os.CreateTemp("", "llcppg.symb.json")
		if err != nil {
			fmt.Printf("Failed to create temp file: %v\n", err)
			return
		}
		tmpFile.Write([]byte(tc.symbFileContent))
		symbolData, err := symbol.GenerateAndUpdateSymbolTable(tc.dylibSymbols, headerSymbolMap, tmpFile.Name())
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println(string(symbolData))
		os.Remove(tmpFile.Name())
	}
}
