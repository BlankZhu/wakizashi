docker build -t probe:1.0 -f ./Dockerfile --build-arg APP_NAME=probe VERSION=1.0 .
docker build -t center:1.0 -f ./Dockerfile --build-arg APP_NAME=center VERSION=1.0 .