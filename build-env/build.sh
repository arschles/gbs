#!/bin/bash

CLONE_URL="https://$SITE/$ORG/$REPO.git"

echo "Cloning $CLONE_URL"
git clone --depth=1 $CLONE_URL
if [[ $? != 0 ]];then
  echo "Error cloning"
  exit 1
fi

mkdir -p /go/src/$SITE/$ORG
mv $REPO /go/src/$SITE/$ORG/$REPO
cd /go/src/$SITE/$ORG/$REPO
echo "Building"

go build -o $BIN_NAME
mv ./$BIN_NAME $BIN_DIR/$BIN_NAME
echo "Done building $SITE/$ORG/$REPO (moved to $BIN_DIR/$BIN_NAME in container)"
