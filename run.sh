set -eu

. `dirname $0`/build.sh $@

APP_NAME=`head -1 go.mod | cut -d ' ' -f 2`

PROJ_DIR=`pwd`

${PROJ_DIR}/${BUILD_DIR}/${APP_NAME}
