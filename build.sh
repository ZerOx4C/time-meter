set -eu

cd `dirname $0`

BUILD_PARAMS=-ldflags\ '-H=windowsgui'
BUILD_DIR=build/release

if [ "${1:-}" = "-d" ]; then
	BUILD_PARAMS=-tags\ debug
	BUILD_DIR=build/debug
fi

git rev-parse --short HEAD | tr -d '\n' > env/revision.txt

mkdir -p ${BUILD_DIR}
go build ${BUILD_PARAMS} -o ${BUILD_DIR} .
