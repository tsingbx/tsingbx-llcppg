package convert

import (
	"math"
	"testing"
)

func TestLitToUint(t *testing.T) {
	type CaseType struct {
		lit  string
		want uint64
		typ  IntType
	}

	int32Min := uint64(1) << 31
	testCases := []CaseType{
		{lit: "123", want: 123, typ: TypeInt},
		{lit: "0xDEEDBEAF", want: 0xDEEDBEAF, typ: TypeUint},                  // DEEDBEAF
		{lit: "0x80000000", want: int32Min, typ: TypeUint},                    // INT32_MIN
		{lit: "0x7FFFFFFF", want: math.MaxInt32, typ: TypeInt},                // INT32_MAX
		{lit: "0xFFFFFFFF", want: math.MaxUint32, typ: TypeUint},              // UINT32_MAX
		{lit: "0xFFFFFFFFFFFFFFFF", want: 0xFFFFFFFFFFFFFFFF, typ: TypeUlong}, // UINT64_MAX
		{lit: "2147483647", want: math.MaxInt32, typ: TypeInt},                // INT32_MAX
		{lit: "9223372036854775807", want: math.MaxInt64, typ: TypeUlong},     // INT64_MAX
	}
	for _, tc := range testCases {
		result, typ, err := litToUint(tc.lit)

		if err != nil {
			t.Error(err)
		}
		if tc.typ != typ {
			t.Error(tc.lit, "type mismatch want", tc.typ, "got", typ)
		}

		if result^(tc.want) != 0 || result != (tc.want) {
			t.Error(tc.lit, "result mismatch want", tc.want, "got", result, typ)
		}
	}
}

func TestLitToInt(t *testing.T) {
	type CaseType struct {
		lit  string
		want int64
		typ  IntType
	}

	testCases := []CaseType{
		{lit: "-2147483648", want: math.MinInt32, typ: TypeInt},           // INT32_MIN
		{lit: "-9223372036854775808", want: math.MinInt64, typ: TypeLong}, // INT64_MIN
	}
	for _, tc := range testCases {
		result, typ, err := litToInt(tc.lit)

		if err != nil {
			t.Error(err)
		}
		if tc.typ != typ {
			t.Error(tc.lit, "type mismatch want", tc.typ, "got", typ)
		}

		if result^(tc.want) != 0 || result != (tc.want) {
			t.Error(tc.lit, "result mismatch want", tc.want, "got", result, typ)
		}
	}
}
