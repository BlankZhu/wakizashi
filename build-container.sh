set -euo pipefail

tag="$1"

if [ -z "$tag" ]; then
    echo "invalid parameters, check again"
    exit 1
fi

docker build -f ./Dockerfile -t wakizashi/probe:${tag}  --build-arg "APP_NAME=probe"  --build-arg "VERSION=${tag}" .
docker build -f ./Dockerfile -t wakizashi/center:${tag} --build-arg "APP_NAME=center" --build-arg "VERSION=${tag}" .