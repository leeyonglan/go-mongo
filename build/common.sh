BIN="teaapp"

TYPE="${1:-dev}"

DIST_DIR="$(dirname "$DIR")/bin"
mkdir -p "$DIST_DIR"

EXE="${DIST_DIR}/${BIN}-${TYPE}"