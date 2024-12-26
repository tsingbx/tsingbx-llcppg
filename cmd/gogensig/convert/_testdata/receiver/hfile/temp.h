struct in_addr1 {
    unsigned int s_addr;
};

struct ares_in6_addr {
    union {
        unsigned char _S6_u8[16];
    } _S6_un;
};

struct ares_addr {
    int family;

    union {
        struct in_addr1 addr4;
        struct ares_in6_addr addr6;
    } addr;
};

#include "use.h"