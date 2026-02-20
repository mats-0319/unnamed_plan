#!/bin/bash
# Use At: Cloud Server
# Behavior: Stop and Restart Servers

pkill -f unnamed_plan_server_exec

from_path=$(pwd)

# change dir, for run this script from anywhere
cd "$(dirname "$0")" || exit 1

  function start_exec() {
    cd "./server" || exit 1
    chmod +x "./unnamed_plan_server_exec"
    (./unnamed_plan_server_exec &) # &: not block. (): still work after cmd exit
    cd .. || exit 1
  }

start_exec

cd "$from_path" || exit 1
