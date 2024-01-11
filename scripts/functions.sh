#!/bin/bash

red=$'\033[1;31m'
green=$'\033[1;32m'
yellow=$'\033[1;33m'
blue=$'\033[1;34m'
off=$'\e[m'

#  Run a command, prefixing all output with a string.
#  Useful for running multiple commands in parallel, and identifying which is which.
#
#      Usage: prefix-with "app1" "green" "npm run start"
#
prefix() { 
    local prefix="$1";
    shift;
    "$@" > >(sed "s/^/$prefix: /") 2> >(sed "s/^/$prefix (err): /" >&2); 
}

# 
# Differetn variants of prefix(), which print the prefix in a specific color
# 
prefix-red() { 
    local prefix="$1";
    shift;
    "$@" > >(sed "s/^/$red$prefix:$off /") 2> >(sed "s/^/$red$prefix (err):$off /" >&2); 
}
prefix-green() { 
    local prefix="$1";
    shift;
    "$@" > >(sed "s/^/$green$prefix:$off /") 2> >(sed "s/^/$green$prefix (err):$off /" >&2); 
}
prefix-blue() { 
    local prefix="$1";
    shift;
    "$@" > >(sed "s/^/$blue$prefix:$off /") 2> >(sed "s/^/$blue$prefix (err):$off /" >&2); 
}
prefix-yellow() { 
    local prefix="$1";
    shift;
    "$@" > >(sed "s/^/$yellow$prefix:$off /") 2> >(sed "s/^/$yellow$prefix (err):$off /" >&2); 
}

#  Find the process which is listening on local TCP port, and kill it
#  Tested on MacOS. 
#
#      Usage: kill-process-by-port 8080
#
kill-process-by-port() {
    netstat -anv | grep $1 | cut -w -f 9 | xargs kill -9
}