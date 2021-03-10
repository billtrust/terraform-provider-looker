#! /bin/sh -eu

# This is needed because we need to know the platform name to know
# where to put the plugin, and the platform names exposed by the go
# runtime are weird and can't be replicated with `uname(1)`.

WORK=$(mktemp -d)
trap "rm -rf \"$WORK\"" EXIT
cd $WORK

cat > platform.go <<EOF
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Printf("%s_%s\n", runtime.GOOS, runtime.GOARCH)
}
EOF

go build -o platform
./platform

