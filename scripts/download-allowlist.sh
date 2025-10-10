#!/bin/sh

# Script to download the domain allowlist during image build
# This script is used in the Dockerfile to fetch the latest allowlist

ALLOWLIST_URL="${ALLOWLIST_URL:-https://mirror.ghproxy.com/https://raw.githubusercontent.com/DaoCloud/public-image-mirror/main/domain.txt}"
OUTPUT_FILE="${OUTPUT_FILE:-/etc/repimage/allowlist.txt}"

echo "Downloading allowlist from: $ALLOWLIST_URL"
echo "Output file: $OUTPUT_FILE"

# Create directory if it doesn't exist
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Try to download the allowlist
if wget -q -O "$OUTPUT_FILE" "$ALLOWLIST_URL"; then
    echo "Successfully downloaded allowlist to $OUTPUT_FILE"
    echo "Total entries: $(grep -c '=' "$OUTPUT_FILE" 2>/dev/null || echo 0)"
else
    echo "Failed to download allowlist from $ALLOWLIST_URL"
    echo "Creating default allowlist..."
    
    # Create a default allowlist if download fails
    cat > "$OUTPUT_FILE" <<EOF
# Default domain allowlist for repimage
docker.io=m.daocloud.io/docker.io
gcr.io=m.daocloud.io/gcr.io
k8s.gcr.io=m.daocloud.io/k8s.gcr.io
registry.k8s.io=m.daocloud.io/registry.k8s.io
ghcr.io=m.daocloud.io/ghcr.io
quay.io=m.daocloud.io/quay.io
EOF
    echo "Created default allowlist with $(grep -c '=' "$OUTPUT_FILE") entries"
fi

exit 0
