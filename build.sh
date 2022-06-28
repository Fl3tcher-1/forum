#!/bin/bash
set -e

# Run this as `sudo ./build.sh` if it gives you permission denied
docker build . -t forum

echo "Container image built, run locally with:
docker run -v $(pwd)/database:/database -p 8080:8080 forum"