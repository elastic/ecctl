#!/bin/bash

echo "-> Creating useful aliases and helper functions for ecctl."

COMPLETION_COMMANDS=( dev-cli ecctl-container )

unset-credentials() {
    unset EC_HOST EC_USER EC_PASS EC_REGION EC_BEARER
}

dev-cli() {
    go run main.go ${*}
}

dump-ecctl-vars() {
    env | grep "EC_"
}

reload_helpers() {
    source "$(git rev-parse --show-toplevel)/scripts/helpers.sh"
}

for command in "${COMPLETION_COMMANDS[@]}"; do
    source <(dev-cli generate completions --binary="${command}")
done

echo "-> Helpers loaded in the environment."
echo "-> To run the latest code, simply type \"dev-cli\"."
