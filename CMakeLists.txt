cmake_minimum_required(VERSION 3.2)

set(CMAKE_MODULE_PATH ${CMAKE_CURRENT_LIST_DIR}/)

project(sgps-core Go)

set(SOURCE_FILES
    src/sgps-core.go)

add_executable(sgps-core ${SOURCE_FILES})
