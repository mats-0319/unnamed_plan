#!/bin/bash
# Use At: Local Dev Env
# Behavior: Generate Server and UI at 'root/build'

from_path=$(pwd)

# change dir, for run this script from anywhere
cd "$(dirname "$0")" || exit 1
cd .. || exit 1

if [ -d "./build/" ]; then
  rm -rf ./build/*
fi

# '-p' flag will generate parent path if not exist
mkdir -p "./build/server/"
# ui will 'mv' whole dir, unnecessary to 'mkdir'

# build server

  cd "./server" || exit 1
  go mod tidy
  cd .. || exit 1

    function build_server() {
      local server_name="$1"

      cd "./server/${server_name}" || exit 1
      go build -o "server_exec"
      cd ../.. || exit 1

      mv "./server/${server_name}/server_exec" "./build/server/unnamed_plan_server_exec"
      cp "./server/${server_name}/config_production.json" "./build/server/config.json"
    }

  build_server "cmd"

# build ui

  cd "./web" || exit 1
  pnpm install
  pnpm run build
  cd .. || exit 1

  mv "./web/dist/" "./build/dist/"

# shell

  cp "./scripts/restart_server.sh" "./build/restart_server.sh"

cd "$from_path" || exit 1
