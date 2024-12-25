typedef unsigned int __uint32_t;
typedef __uint32_t in_addr_t;
struct in_addr {
    in_addr_t s_addr;
};

struct ares_in6_addr {
    union {
        unsigned char _S6_u8[16];
    } _S6_un;
};

struct ares_addr {
    int family;

    union {
        struct in_addr addr4;
        struct ares_in6_addr addr6;
    } addr;
};

#include "use.h"
