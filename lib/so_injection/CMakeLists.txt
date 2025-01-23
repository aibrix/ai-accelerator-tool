# Add minimum CMake version
cmake_minimum_required(VERSION 3.10)

# Add project name
project(gpu-diagnosis-injection)

# Fix policy warning for boolean values
cmake_policy(SET CMP0012 NEW)

set(CXX_OPT 3)
# TODO now we link nvml/dcmi/... at link time thus them must be present if we set the coresponding option to on
option(TEST_UTILS_ENABLE_NVIDIA "Enable Nvidia device test" OFF)

option(TEST_UTILS_INJECT_DL "Enable inject dl" ON)

get_filename_component(GPU_DIAG_JOB_ROOT_DIR ${CMAKE_CURRENT_SOURCE_DIR} PATH)

set(TOMLPLUSPLUS_INCLUDE_PATH ${GPU_DIAG_JOB_ROOT_DIR}/third_party/marzer/tomlplusplus/include)
set(NVDIAML_INCLUDE_PATH ${GPU_DIAG_JOB_ROOT_DIR}/third_party/device_mgr/nvml/include)
set(NVDIAML_LIB_PATH ${GPU_DIAG_JOB_ROOT_DIR}/third_party/device_mgr/nvml/lib)

set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${GPU_DIAG_JOB_ROOT_DIR}/build/lib)

# Add C++11 support
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

add_library(devso-injection SHARED
    ${CMAKE_CURRENT_SOURCE_DIR}/inject_helper.cpp
)

target_include_directories(devso-injection PUBLIC ${GPU_DIAG_JOB_ROOT_DIR}/test_utils/so_injection/common)
target_include_directories(devso-injection PUBLIC ${TOMLPLUSPLUS_INCLUDE_PATH})
set_property(TARGET devso-injection PROPERTY POSITION_INDEPENDENT_CODE ON)

target_compile_options(devso-injection PRIVATE -Wall -O${CXX_OPT})
target_link_options(devso-injection PRIVATE -Wall -nostartfiles -O${CXX_OPT})

if(${TEST_UTILS_ENABLE_NVIDIA})
    if(APPLE)
        message(WARNING "NVIDIA support disabled on macOS")
    else()
        # Check if NVIDIA library exists
        if(EXISTS "${NVDIAML_LIB_PATH}/libnvidia-ml.so")
            target_compile_definitions(devso-injection PUBLIC INJECT_NVML)
            target_include_directories(devso-injection PUBLIC ${NVDIAML_INCLUDE_PATH})
            target_link_libraries(devso-injection ${NVDIAML_LIB_PATH}/libnvidia-ml.so)
        else()
            message(WARNING "NVIDIA library not found at ${NVDIAML_LIB_PATH}/libnvidia-ml.so - building without NVIDIA support")
            set(TEST_UTILS_ENABLE_NVIDIA OFF)
        endif()
    endif()
endif()

if(${TEST_UTILS_INJECT_DL})
    target_link_libraries(devso-injection ${CMAKE_DL_LIBS})
    target_compile_definitions(devso-injection PUBLIC INJECT_DL)
endif()

add_custom_target(
    copy_lib_to_cli_resource ALL
    COMMAND ${CMAKE_COMMAND} -E copy 
        $<TARGET_FILE:devso-injection>
        ${GPU_DIAG_JOB_ROOT_DIR}/test_utils/cli/resource/$<TARGET_FILE_NAME:devso-injection>
    DEPENDS devso-injection 
)