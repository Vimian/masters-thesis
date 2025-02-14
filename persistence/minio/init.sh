#!/bin/bash
# set minio environment variables
MINIO_ROOT_USER=${MINIO_ROOT_USER} #
MINIO_ROOT_PASSWORD=${MINIO_ROOT_PASSWORD} #
MINIO_SERVER=${MINIO_SERVER} #
# set minio alias
mc alias set minio http://$MINIO_SERVER $MINIO_ROOT_USER $MINIO_ROOT_PASSWORD #
# create bucket
BUCKET_NAME="benchmark" #
mc mb minio/"$BUCKET_NAME" #
# upload original files
mc cp /minio/original minio/"$BUCKET_NAME" --recursive #