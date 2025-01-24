set(CXX_OPT 3)
option(TEST_UTILS_ENABLE_NVIDIA "Enable Nvidia device test" ON)
option(TEST_UTILS_ENABLE_ASCEND "Enable Ascend device test" OFF)

option(TEST_UTILS_INJECT_DL "Enable inject dl" ON)

set(TOMLPLUSPLUS_INCLUDE_PATH ${PROJECT_ROOT_DIR}/lib/third_party/marzer/tomlplusplus/include)
set(NVDIAML_INCLUDE_PATH ${PROJECT_ROOT_DIR}/lib/third_party/device_mgr/nvml/include)
set(NVDIAML_LIB_PATH ${PROJECT_ROOT_DIR}/lib/third_party/device_mgr/nvml/lib)
set(DCMI_INCLUDE_PATH ${PROJECT_ROOT_DIR}/lib/third_party/device_mgr/dcmi/include)
set(DCMI_LIB_PATH ${PROJECT_ROOT_DIR}/lib/third_party/device_mgr/dcmi/lib)

set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${PROJECT_ROOT_DIR}/lib/build/lib)


add_library(devso-injection SHARED ${PROJECT_ROOT_DIR}/lib/src/dev_mgr_hijack/nvml/injection.cpp ${PROJECT_ROOT_DIR}/lib/src/dev_mgr_hijack/dcmi/injection.cpp ${PROJECT_ROOT_DIR}/lib/src/dev_mgr_hijack/inject_helper.cpp)

target_include_directories(devso-injection PUBLIC ${PROJECT_ROOT_DIR}/lib/src/dev_mgr_hijack/common)
target_include_directories(devso-injection PUBLIC ${TOMLPLUSPLUS_INCLUDE_PATH})
set_property(TARGET devso-injection PROPERTY POSITION_INDEPENDENT_CODE ON)

target_compile_options(devso-injection PRIVATE -Wall -O${CXX_OPT})
target_link_options(devso-injection PRIVATE -Wall -nostartfiles -O${CXX_OPT})

if(${TEST_UTILS_ENABLE_NVIDIA})
    target_compile_definitions(devso-injection PUBLIC INJECT_NVML)
    target_include_directories(devso-injection PUBLIC ${NVDIAML_INCLUDE_PATH})
    target_link_libraries(devso-injection ${NVDIAML_LIB_PATH}/libnvidia-ml.so)
endif()

if (${TEST_UTILS_ENABLE_ASCEND})
    target_compile_definitions(devso-injection PUBLIC INJECT_DCMI)
    target_include_directories(devso-injection PUBLIC ${DCMI_INCLUDE_PATH})
    target_link_libraries(devso-injection ${DCMI_LIB_PATH}/libdcmi.so)
endif()

if(${TEST_UTILS_INJECT_DL})
    target_link_libraries(devso-injection ${CMAKE_DL_LIBS})
    target_compile_definitions(devso-injection PUBLIC INJECT_DL)
endif()

add_custom_target(
    copy_lib_to_cli_resource ALL
    COMMAND cp ${CMAKE_LIBRARY_OUTPUT_DIRECTORY}/libdevso-injection.so ${PROJECT_ROOT_DIR}/pkg/mock/resources/libdevso-injection.so 
    DEPENDS devso-injection 
)