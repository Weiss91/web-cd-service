#!/bin/bash
# update build files
bazelisk run //:gazelle

# update build deps
bazelisk run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

# set needed envs to start
export GIT_USER="user"
export GIT_PASSWORD="pw"
export GIT_DATA_PATH="sample/path"
# --> the path where the git repos should be stored that bazel can execute it.
# Defaults to /git)

read -d '' conf << EOF
{
    "path": "path",
    "registries": {
            "develop": {
                    "Registry": "develop",
                    "user": "devRegistryUser",
                    "password": "devRegistryPassword"
            },
            "production": {
                    "Registry": "production",
                    "user": "prodRegistryUser",
                    "password": "prodRegistryPassword"
            }
    }
}
EOF
export DOCKER_CONF=$conf
export API_KEY_READ="READAPIKEY"
export API_KEY_EXEC="EXECAPIKEY"
export SERVER_PORT="8088"
# --> the Port on which the web-api is running
export STORAGE_PATH= "/data"
# --> the directory where the active and historic tasks are stored

bazelisk.exe run bazelwebapi:bazelwebapi
