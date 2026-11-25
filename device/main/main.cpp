// DRAK device entry point. S0 skeleton: no CommandCore, no IR yet.
#include <cstdio>

#include "drak/id.hpp"

extern "C" void app_main(void) {
    const char* device_id = "room-01";
    std::printf("drak-device skeleton: device_id=%s valid=%d\n", device_id,
                drak::is_valid_id(device_id) ? 1 : 0);
}
