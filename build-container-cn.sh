set -euo pipefail

tag="$1"

if [ -z "$tag" ]; then
    echo "invalid parameters, check again"
    exit 1
fi

docker build -f ./Dockerfile.cn -t wakizashi/probe:${tag}  --build-arg "APP_NAME=probe"  --build-arg "VERSION=${tag}" --network host .
docker build -f ./Dockerfile.cn -t wakizashi/center:${tag} --build-arg "APP_NAME=center" --build-arg "VERSION=${tag}" --network host .