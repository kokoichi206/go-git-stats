#!/bin/bash
#
# Description:
#   Add a new tag and push it to origin
#
# Usage:
#   bash tag_push.sh 0.1.4
#
set -euo pipefail

# No arguments
if [[ "$#" = 0 ]]; then
    echo "You need to give 1 argument as a new tag version."
    echo "e.g.) bash tag_push.sh 0.1.4"
    exit 1
fi

# Wrong arguments
if [[ ! "$1" =~ "v"?([0-9]\.[0-9]\.[0-9]) ]]; then
    echo "new tag version $1 is not correct."
    echo "please follow the semantic versioning rule (e.g. 0.1.4 or v0.1.4)"
    exit 1
fi

# add tag
git tag "v${BASH_REMATCH[1]}"
# push tag
git push origin --tag
