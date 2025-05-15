#include "type.h"
#include <cJSON.h>
#include <stddef.h>
#include <thirddep.h>
#include <thirddep2.h>
#include <thirddep3.h>
// This file is supposed to depend on cjson in its cflags, but for testing,
// we will simulate its API using libcjson instead.
//   "cflags" :"$(pkg-config --cflags libcjson)"
cJSON *create_response(int status_code, const char *message);

cJSON_bool parse_client_request(const char *json_string, char *error_buffer, size_t buffer_size);

cJSON_bool serialize_response(cJSON *response, char *buffer, const int length, const cJSON_bool pretty_print);

third_dep third_depfn(third_dep *a, third_dep2 *b, _depcjson_type c, basic_dep d);

third_dep3 third_type(third_dep3 *a);
    
// This struct demonstrates the handling of same llcppg.pub names across different packages:
//
// 1. Basic_stream (from basicdep.h)
//    - Indirect dependency
//    - llcppg.pub mapping: Basic_stream -> Stream
//
// 2. third_dep_stream (from thirddep.h)
//    - Direct dependency
//    - llcppg.pub mapping: third_dep_stream -> Stream
typedef struct samePubStream {
    Basic_stream basic_stream;
    third_dep_stream third_dep_stream;
} samePubStream;

