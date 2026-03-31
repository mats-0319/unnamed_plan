#!/bin/bash
# Use At: Local Dev Env
# Behavior: Compile Ebiten Game and 'cp' to its Position

from_path=$(pwd) # record current path

# change dir, for run this script from anywhere
cd "$(dirname "$0")" || exit 1
cd .. || exit 1

# compile game

  function compile_game() {
    local game_name="$1"

    cd "./game/${game_name}" || exit 1
    go mod tidy
    env GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o "${game_name}.wasm"
    cd ../.. || exit

    if [ -f "./web/public/${game_name}.wasm" ]; then
      rm "./web/public/${game_name}.wasm"
    fi
    mv "./game/${game_name}/${game_name}.wasm" "./web/public/${game_name}.wasm"

    if [ -f "./web/public/${game_name}.html" ]; then
      rm "./web/public/${game_name}.html"
    fi
    cp "./game/${game_name}/${game_name}.html" "./web/public/${game_name}.html"

    if [ ! -f "./web/public/wasm_exec.js" ]; then
      cp "./game/wasm_exec.js" "./web/public/wasm_exec.js"
    fi
  }

compile_game "flip"

# back to from path
cd "$from_path" || exit 1
