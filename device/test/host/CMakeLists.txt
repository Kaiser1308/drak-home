# Host (PC) unit tests for device/components/*, built without ESP-IDF.
#   cmake -S device/test/host -B device/build-host
#   cmake --build device/build-host
#   ctest --test-dir device/build-host --output-junit "$PWD/device-host.xml"
cmake_minimum_required(VERSION 3.21)
project(drak_device_host_tests CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
FetchContent_Declare(
  googletest
  URL https://github.com/google/googletest/releases/download/v1.17.0/googletest-1.17.0.tar.gz
  URL_HASH SHA256=65fab701d9829d38cb77c14acdc431d2108bfdbf8979e40eb8ae567edf10b27c
  DOWNLOAD_EXTRACT_TIMESTAMP TRUE)
set(INSTALL_GTEST OFF CACHE BOOL "" FORCE)
set(BUILD_GMOCK OFF CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

set(DRAK_COMPONENTS ${CMAKE_CURRENT_SOURCE_DIR}/../../components)

add_library(drak_core STATIC ${DRAK_COMPONENTS}/drak_core/src/id.cpp)
target_include_directories(drak_core PUBLIC ${DRAK_COMPONENTS}/drak_core/include)
target_compile_options(drak_core PRIVATE -Wall -Wextra -Werror)

enable_testing()
add_executable(drak_core_tests id_test.cpp)
target_link_libraries(drak_core_tests PRIVATE drak_core GTest::gtest_main)
target_compile_options(drak_core_tests PRIVATE -Wall -Wextra -Werror)

include(GoogleTest)
gtest_discover_tests(drak_core_tests)
