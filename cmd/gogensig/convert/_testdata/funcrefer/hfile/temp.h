#include <stddef.h>
#include <stdio.h>
typedef int (*CallBack)(void *L);
void exec(void *L, CallBack cb);

typedef struct Hooks
{
  void *(*malloc_fn)(size_t sz);
  void (*free_fn)(void *ptr);
} Hooks;

typedef struct Stream
{
  FILE *f;
  CallBack cb;
} Stream;

typedef void(OSSL_provider_init_fn)();
extern OSSL_provider_init_fn OSSL_provider_init2;

typedef struct OSSL_CORE_HANDLE OSSL_CORE_HANDLE;
typedef struct OSSL_DISPATCH OSSL_DISPATCH;
typedef int(OSSL_provider_init_fn2)(const OSSL_CORE_HANDLE *handle,
                                   const OSSL_DISPATCH *in,
                                   const OSSL_DISPATCH **out,
                                   void **provctx);

OSSL_provider_init_fn2 OSSL_provider_init;