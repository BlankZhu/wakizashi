set - euo pipefail

version="$1"

if [ -z "$version" ]; then
    echo "invalid parameters, check again"
    exit 1
fi

echo "generating wakizashi center"
GO111MODULE=on go build -ldflags \
    "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=${version} -X main.gitCommitID=`git rev-parse HEAD`" \
    -o ./build/center \
    ./cmd/center
echo "wakizashi center generated"

echo "generating wakizashi probe"
GO111MODULE=on go build -ldflags \
    "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=${version} -X main.gitCommitID=`git rev-parse HEAD`" \
    -o ./build/probe \
    ./cmd/probe
echo "wakizashi probe generated"