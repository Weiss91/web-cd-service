#!/bin/bash
# update build files
bazelisk run //:gazelle

# update build deps
bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

# set needed envs to start
export GIT_USER="user"
export GIT_PASSWORD="pw"
export GIT_DATA_PATH="sample/path"
export DOCKER_CONF_PATH="./config.json"
export API_KEY_READ="READAPIKEY"
export API_KEY_EXEC="EXECAPIKEY"
export SERVER_PORT="8088"

bazelisk.exe run bazelwebapi:bazelwebapi
