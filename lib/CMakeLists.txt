cmake_minimum_required(VERSION 3.18 FATAL_ERROR)
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(CMAKE_C_COMPILER "/usr/bin/clang")
set(CMAKE_CXX_COMPILER "/usr/bin/clang++")

project(gpu_diagnosis_jobs)

if(CMAKE_SYSTEM_PROCESSOR MATCHES "x86_64")
  set(TARGET_ARCH "x86_64-linux-gnu")
elseif(CMAKE_SYSTEM_PROCESSOR MATCHES "aarch64")
  set(TARGET_ARCH "aarch64-linux-gnu")
endif()

get_filename_component(PROJECT_ROOT_DIR ${CMAKE_CURRENT_SOURCE_DIR} PATH)

add_subdirectory(${PROJECT_ROOT_DIR}/lib/src/dev_mgr_hijack ${PROJECT_ROOT_DIR}/lib/build/dev_mgr_hijack)