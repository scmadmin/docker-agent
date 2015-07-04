PROJECT( docker-agent )
cmake_minimum_required(VERSION 3.0)

# Set project layout variables
SET(BII_PROJECT_ROOT "/home/paul/dev/docker-agent")
SET(BII_BLOCKS_DIR blocks)
SET(BII_DEPS_DIR deps)
SET(BII_LIB_DIR lib)
SET(BII_CMAKE_DIR )
SET(BII_BIN_DIR bin)
set(CMAKE_MODULE_PATH "${BII_PROJECT_ROOT}/blocks"
					  "${BII_PROJECT_ROOT}/deps")

# inclusion of general biicode macros, biicode.cmake 
INCLUDE(biicode.cmake) 
IF (${CMAKE_SYSTEM_NAME} STREQUAL "Arduino")
	OPTION(BII_DEPS_SYSTEM "Handle deps as SYSTEM, avoiding warnings" OFF)
ELSE()
	OPTION(BII_DEPS_SYSTEM "Handle deps as SYSTEM, avoiding warnings" ON)
ENDIF()
SET(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${BII_PROJECT_ROOT}/bin)
SET(CMAKE_RUNTIME_OUTPUT_DIRECTORY_RELEASE ${BII_PROJECT_ROOT}/bin)
SET(CMAKE_RUNTIME_OUTPUT_DIRECTORY_DEBUG ${BII_PROJECT_ROOT}/bin)

SET(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${BII_PROJECT_ROOT}/lib)
SET(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_RELEASE ${BII_PROJECT_ROOT}/lib)
SET(CMAKE_ARCHIVE_OUTPUT_DIRECTORY_DEBUG ${BII_PROJECT_ROOT}/lib)

# Including /blocks and /deps directories
INCLUDE_DIRECTORIES(${BII_PROJECT_ROOT}/blocks)
IF(BII_DEPS_SYSTEM)
	INCLUDE_DIRECTORIES(SYSTEM ${BII_PROJECT_ROOT}/deps)
ELSE()
	INCLUDE_DIRECTORIES(${BII_PROJECT_ROOT}/deps)
ENDIF()
ADD_DEFINITIONS(-DBIICODE)
SET(BIICODE_ENV_DIR "/home/paul/.biicode")

#artifact to allow upstream configurations
BII_PREBUILD_STEP(deps/hjiang/jsonxx)
BII_PREBUILD_STEP(deps/lasote/simple_logger)
BII_PREBUILD_STEP(deps/lasote/stringutils)
BII_PREBUILD_STEP(deps/lasote/http_parser)
BII_PREBUILD_STEP(deps/lasote/libuv)
BII_PREBUILD_STEP(deps/lasote/httpmodels)
BII_PREBUILD_STEP(deps/lasote/async_http_request)
BII_PREBUILD_STEP(deps/lasote/lambda_http_client)
BII_PREBUILD_STEP(deps/lasote/docker_client)
BII_PREBUILD_STEP(blocks/devadaptive/docker-agent)

# Enable tests with CTest
enable_testing()
# Creating a custom target which will bring associated
# all the test targets
add_custom_target(biitest)

# Inclusion of the blocks affected by the build
BII_INCLUDE_BLOCK(deps/hjiang/jsonxx)
BII_INCLUDE_BLOCK(deps/lasote/simple_logger)
BII_INCLUDE_BLOCK(deps/lasote/stringutils)
BII_INCLUDE_BLOCK(deps/lasote/http_parser)
BII_INCLUDE_BLOCK(deps/lasote/libuv)
BII_INCLUDE_BLOCK(deps/lasote/httpmodels)
BII_INCLUDE_BLOCK(deps/lasote/async_http_request)
BII_INCLUDE_BLOCK(deps/lasote/lambda_http_client)
BII_INCLUDE_BLOCK(deps/lasote/docker_client)
BII_INCLUDE_BLOCK(blocks/devadaptive/docker-agent)
