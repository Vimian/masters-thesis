#!/bin/bash

# set minio environment variables
MINIO_ROOT_USER=${MINIO_ROOT_USER}
MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD}
MINIO_SERVER=${MINIO_SERVER:-minio:9000}

# wait for minio to be ready
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

# set minio alias
mc alias set minio http://$MINIO_SERVER $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD

# create bucket
BUCKET_NAME="benchmark"
mc mb minio/"$BUCKET_NAME"

# upload original files
mc cp /minio/original minio/"$BUCKET_NAME" --recursive