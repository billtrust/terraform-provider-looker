#!/bin/bash
set -e

if [ -z "$1" ]
then
    echo "'Version' must be specified. No further action taken."
    echo ""
    echo "USAGE: ./build.sh <version>"
    echo ""
    echo "EXAMPLE:"
    echo "   ./build.sh 9.9.9"
    echo ""
    exit -1
else
    VERSION=$1
    echo "Version specified: $VERSION"
fi


echo "== Building Docker Image"
docker build -t terraform-provider-looker -f Dockerfile .

echo "== Building: linux"
docker run -it -v $(pwd):/go/src/github.com/billtrust/terraform-provider-looker terraform-provider-looker bash -c "go get && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/terraform-provider-looker_v$VERSION"

echo "== Building: darwin"
docker run -it -v $(pwd):/go/src/github.com/billtrust/terraform-provider-looker terraform-provider-looker bash -c "go get && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/darwin/terraform-provider-looker_v$VERSION"

echo "== Building: windows"
docker run -it -v $(pwd):/go/src/github.com/billtrust/terraform-provider-looker terraform-provider-looker bash -c "go get && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-looker_v$VERSION.exe"

# Workaround for Docker, sometime resulting binaries are only owned by root.
# This ensures that ownership is reassigned correctly to the user. 
docker run -it -v $(pwd):/go/src/github.com/billtrust/terraform-provider-looker terraform-provider-looker bash -c "chmod -R a+rwX bin/"