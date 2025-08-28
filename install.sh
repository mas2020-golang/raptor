#!/usr/bin/env bash

#########################
# Repo specific content #
#########################

export ALIAS_NAME="crptx"
export OWNER=mas2020-golang
export REPO=raptor
export BINLOCATION="/usr/local/bin"
export APP_NAME="raptor"
export SUCCESS_CMD="${BINLOCATION}/${APP_NAME} version"

# -- COLORS
export STOP_COLOR="\e[0m"
# color for a main activity
export ACTIVITY="\e[38;5;184m"
# color for a sub activity
export SUB_ACT="\n\n\e[1;34m➜\e[0m"
export DONE="\e[1;32m✔︎\e[0m"
export ERROR="\e[1;31mERROR\e[0m:"
export WARNING="\e[38;5;216mWARNING\e[0m:"

###############################
# Content common across repos #
###############################set -x
printf "${ACTIVITY}%s ${STOP_COLOR}" "installation for the $REPO application..."
version=$(curl -sI https://github.com/$OWNER/$REPO/releases/latest | grep -i "location:" | awk -F"/" '{ printf "%s", $NF }' | tr -d '\r')
#set -x
printf "\nselected version for %s is '%q'" $REPO $version
if [ ! $version ]; then
    echo "Failed while attempting to install $REPO. Please manually install:"
    echo ""
    echo "1. Open your web browser and go to https://github.com/$OWNER/$REPO/releases"
    echo "2. Download the latest release for your platform. Call it '$REPO'."
    echo "3. chmod +x ./$REPO"
    echo "4. mv ./$REPO $BINLOCATION"
    if [ -n "$ALIAS_NAME" ]; then
        echo "5. ln -sf $BINLOCATION/$REPO /usr/local/bin/$ALIAS_NAME"
    fi
    exit 1
fi

hasCli() {
    hasCurl=$(which curl)
    if [ "$?" = "1" ]; then
        echo "You need curl to use this script."
        exit 1
    fi
}

getPackage() {
    uname=$(uname)
    userid=$(id -u)

    suffix=""
    case $uname in
    "Darwin")
        arch=$(uname -m)
        case $arch in
        "x86_64")
            suffix="Darwin_x86_64"
            ;;
        esac
        case $arch in
        "arm64")
            suffix="Darwin_arm64"
            ;;
        esac
        ;;

    "MINGW"*)
        suffix=".exe"
        BINLOCATION="$HOME/bin"
        mkdir -p $BINLOCATION

        ;;
    "Linux")
        arch=$(uname -m)
        case $arch in
        "aarch64")
            suffix="Linux_arm64"
            ;;
        esac
        case $arch in
        "x86_64")
            suffix="Linux_x86_64"
            ;;
        esac
        ;;
    esac

    #cryptex_0.1.0-rc.1_Linux-x86_64.tar.gz
    downloadFile="${REPO}_${suffix}.zip"
    targetFile="/tmp/${downloadFile}"
    printf "\nthe file to download is '%q'" "${downloadFile}"

    url="https://github.com/$OWNER/$REPO/releases/download/$version/${downloadFile}"
    printf "${SUB_ACT} %s ${STOP_COLOR}" "downloading package $url as ${targetFile}..."
    http_code=$(curl -sSL $url -w '%{http_code}\n' --output "${targetFile}")

    # check the file not found
    if [ "$?" != "0" ] || [ ${http_code} -eq 404 ]; then
        printf "\n${ERROR} no file as a target download has been found"
        exit 1
    fi

    printf "\n${DONE} download complete\n"

    # unzip the file
    cd /tmp
    unzip "${targetFile}"
    if [ "$?" != "0" ]; then
        printf "\n${ERROR} unzip file"
        exit 1
    fi
    chmod +x "/tmp/${APP_NAME}"
    rm "${targetFile}"

    if [ ! -w "$BINLOCATION" ]; then
        echo
        echo "============================================================"
        echo "  The script was run as a user who is unable to write"
        echo "  to $BINLOCATION. To complete the installation the"
        echo "  following commands may need to be run manually."
        echo "============================================================"
        echo
        printf "${SUB_ACT} %s ${STOP_COLOR}" "moving ${APP_NAME} to $BINLOCATION..."
        
        sudo mv "/tmp/${APP_NAME}" "${BINLOCATION}"

        if [ -e "${APP_NAME}" ]; then
            rm "/tmp/${APP_NAME}"
        fi

        if [ -n "$ALIAS_NAME" ]; then
            echo "  sudo ln -sf ${BINLOCATION}/${APP_NAME} ${BINLOCATION}/${ALIAS_NAME}"
        fi
    else
        printf "${SUB_ACT} %s ${STOP_COLOR}" "moving ${APP_NAME} to ${BINLOCATION}..."

        if [ ! -w "${BINLOCATION}" ] && [ -f "${BINLOCATION}/${APP_NAME}" ]; then
            echo
            echo "================================================================"
            echo "  $${BINLOCATION}/${APP_NAME} already exists and is not writeable"
            echo "  by the current user.  Please adjust the binary ownership"
            echo "  or run sh/bash with sudo."
            echo "================================================================"
            echo
            exit 1
        fi

        cp "/tmp/${APP_NAME}" "${BINLOCATION}/"
        if [ "$?" = "0" ]; then
            printf "\n${DONE} new version of "${APP_NAME}" installed to ${BINLOCATION}"
        fi

        if [ -e "${APP_NAME}" ]; then
            rm "/tmp/${APP_NAME}"
        fi
        
        printf "${SUB_ACT} checking application...\n"
        ${SUCCESS_CMD}
        if [ "$?" != "0" ]; then
            printf "\n${ERROR} the application is not correctly installed"
        else
            printf "\n${DONE} %s successfully installed" "${APP_NAME}"
        fi
    fi
}

hasCli
getPackage
