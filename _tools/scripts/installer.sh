#!/bin/sh
#
# Description:
#  Download binary file from Github release page
#  and move it to the local bin place. (check /usr/local/bin)
#
# Usage:
#  bash installer
#
set -uo pipefail

INSTALL_TO="/usr/local/bin"
WORKING_DIR="${HOME}/.ggs"

function exit_handler() {
    EXIT_CODE="$?"
    if [[ "$EXIT_CODE" -gt 0 ]]; then
        echo "Something went wrong."
    else
        echo "Installed correctly."
    fi
    if [ -d "${WORKING_DIR}" ]; then
        rm -r "${WORKING_DIR}"
    fi
}
trap exit_handler EXIT

RELEASE_URL="https://github.com/kokoichi206/go-git-stats/releases"
# /latest redirects to the newest tag url
## VERSION format is x.x.x
VERSION="$(curl -sfL -o /dev/null -w %{url_effective} ${RELEASE_URL}/latest |\
    grep -o '[0-9]\.[0-9]\.[0-9]')"
exit_status=$?
if [ "${exit_status}" -ne 0 ]; then
    echo "Failed to get the latest version."
    echo "Please check ${RELEASE_URL}."
    exit 1
fi

echo "Downloading ggs v${VERSION} ..."

mkdir "${WORKING_DIR}"
exit_status=$?
if [ "${exit_status}" -ne 0 ]; then
    echo "The folder ${WORKING_DIR} already exists."
    echo "Please remove it and try again."
    exit 1
fi

# download checksums file
curl -sfLo "${WORKING_DIR}/checksums.txt" "${RELEASE_URL}/download/v${VERSION}/checksums.txt"

# check the architecture.
arch="$(uname -m)"
if [ "${arch,,}" = "aarch64" ]; then
    arch="arm64"
fi
# download tar.gz file
TAR_FILE="${WORKING_DIR}/go-git-stats_${VERSION}_$(uname -s)_${arch}.tar.gz"
TAR_URL="${RELEASE_URL}/download/v${VERSION}/go-git-stats_${VERSION}_$(uname -s)_${arch}.tar.gz"
curl -sfLo "${TAR_FILE}" "${TAR_URL}"

# check sha value
## FIXME: using cd is not the best option. (but for checksum)
cd "${WORKING_DIR}"
sha256sum --ignore-missing --quiet --check "${WORKING_DIR}/checksums.txt"
exit_status=$?
cd - > /dev/null
# If the exit status is greater than 0, it means file doesn't exist or checksum is wrong.
if [ "${exit_status}" -ne 0 ]; then
    echo "Something went wrong when downloading."
    exit 1
fi
tar -xf "${TAR_FILE}" -C "${WORKING_DIR}"

echo "ggs command was installed to ${INSTALL_TO}."
sudo mv "${WORKING_DIR}/ggs" "${INSTALL_TO}"
