#!/bin/bash

# Set Minio environment variables
MINIO_ROOT_USER=${MINIO_ROOT_USER}
MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD}
MINIO_SERVER=${MINIO_SERVER:-minio:9000}

# Wait for Minio to be ready
while true; do
  mc ls minio > /dev/null 2>&1
  if [[ $? -eq 0 ]]; then
    echo "mc is ready."
    break
  else
    echo "mc is not ready yet!"
    sleep 5
  fi
done

sleep 5

# Set Minio alias
mc alias set minio http://$MINIO_SERVER $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD

# Create bucket
BUCKET_NAME="benchmark"
mc mb minio/"$BUCKET_NAME"

# Upload original files
mc cp /minio/original minio/"$BUCKET_NAME" --recursive
