#include <stddef.h>
#include <stdio.h>
typedef int (*CallBack)(void *L);
void exec(void *L, CallBack cb);

typedef struct Hooks
{
      void *( *malloc_fn)(size_t sz);
      void (*free_fn)(void *ptr);
}Hooks;

typedef struct Stream {
  FILE *f; 
  CallBack cb; 
} Stream;