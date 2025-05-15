#include <stddef.h>
struct struct1
{
    char *b;
    size_t n;
    union
    {
        long l;
        char b[60];
    } init;
};

struct struct2
{
    char *b;
    size_t size;
    size_t n;
    struct
    {
        long l;
        char b[60];
        struct1 rec;
    } init;
};

union union1
{
    char *b;
    size_t size;
    size_t n;
    struct
    {
        long l;
        char b[60];
        struct2 rec;
    } init;
};

union union2
{
    char *b;
    size_t size;
    size_t n;
    union
    {
        long l;
        char b[60];
        struct2 rec;
    } init;
};