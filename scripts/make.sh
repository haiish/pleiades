#!/bin/bash

# ✨ Always start from project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PROTO_DIR="$PROJECT_ROOT/proto"
OUT_DIR="$PROJECT_ROOT/gen"

# ✨ Create output directory if not exist
mkdir -p "$OUT_DIR"

# ✨ Move into proto folder
cd "$PROTO_DIR"

# ✨ Find and compile ALL .proto files
for proto_file in *.proto; do
    protoc --proto_path="$PROTO_DIR" --go_out="$OUT_DIR" --go-grpc_out="$OUT_DIR" "$proto_file"
    echo "Generated ✨ $proto_file ✨"
done

echo "🌸 All proto files generated successfully into $OUT_DIR! 🍡"
