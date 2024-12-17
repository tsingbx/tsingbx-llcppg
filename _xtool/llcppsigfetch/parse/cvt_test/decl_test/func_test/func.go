package main

import (
	test "github.com/goplus/llcppg/_xtool/llcppsigfetch/parse/cvt_test"
)

func main() {
	TestFuncDecl()
}

func TestFuncDecl() {
	testCases := []string{
		`void foo();`,
		`void foo(int a);`,
		`void foo(int a,...);`,
		`float* foo(int a,double b);`,
		`static inline int add(int a, int b);`,
		`typedef void (fntype)();
		 fntype foo;
		`,
		`typedef long (fntype)(long a);
		 typedef fntype fntype2;
		 fntype2 foo;
	   `,
		`
		typedef struct OSSL_CORE_HANDLE OSSL_CORE_HANDLE;
		typedef struct OSSL_DISPATCH OSSL_DISPATCH;
		typedef int (OSSL_provider_init_fn)(const OSSL_CORE_HANDLE *handle,
		                                const OSSL_DISPATCH *in,
		                                const OSSL_DISPATCH **out,
		                                void **provctx);
		OSSL_provider_init_fn OSSL_provider_init;
		   `,
	}
	test.RunTest("TestFuncDecl", testCases)
}
