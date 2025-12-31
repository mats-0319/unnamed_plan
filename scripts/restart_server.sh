#!/bin/bash
# Use At: Cloud Server
# Behavior: Stop and Restart Servers

pkill unnamed_plan_server_exec

from_path=$(pwd)

# change dir, for run this script from anywhere
cd "$(dirname "$0")" || exit 1

  function start_exec() {
    local server_name="$1"

    cd "./server_${server_name}" || exit 1
    chmod +x "./unnamed_plan_server_exec"
    (./unnamed_plan_server_exec &)
    cd .. || exit 1
  }

start_exec "gateway"
start_exec "1_user"
start_exec "2_note"

cd "$from_path" || exit 1
