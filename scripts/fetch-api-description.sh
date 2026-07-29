#!/usr/bin/env bash
#
# Fetch GL.iNet's machine-readable 4.x API description.
#
# GL.iNet's official reference at dev.gl-inet.com/router-4.x-api is gone. The
# description survives inside python-glinet, which exports GL.iNet's own
# documentation database.
#
# The file is written to /tmp rather than into this repository on purpose:
# python-glinet is GPL-3.0 and gogl is MIT, so vendoring it would raise a
# licensing question. docs/api/ is generated from it instead.
set -euo pipefail

URL="https://raw.githubusercontent.com/tomtana/python-glinet/main/pyglinet/api/api_description.json"
DEST="${1:-/tmp/gl-api-description.json}"

echo "fetching GL.iNet 4.x API description"
echo "  from: $URL"
echo "  to:   $DEST"

curl -fsSL "$URL" -o "$DEST"

size=$(wc -c < "$DEST")
if [ "$size" -lt 100000 ]; then
	echo "error: got only $size bytes; expected roughly 1.4 MB" >&2
	exit 1
fi

groups=$(python3 -c 'import json,sys; print(len(json.load(open(sys.argv[1]))))' "$DEST")
methods=$(python3 -c '
import json, sys
d = json.load(open(sys.argv[1]))
print(sum(len(g.get("case_groups_data") or {}) for g in d.values()))
' "$DEST")

echo "ok: $size bytes, $groups groups, $methods methods"
echo
echo "regenerate the reference with:"
echo "  python3 scripts/generate-api-docs.py $DEST"
