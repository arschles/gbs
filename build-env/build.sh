#!/bin/bash

CLONE_URL="https://$SITE/$ORG/$REPO.git"

echo "Cloning $CLONE_URL"
git clone $CLONE_URL
if [[ $? != 0 ]];then
  echo "Error cloning"
  exit 1
fi

mkdir -p /go/src/$SITE/$ORG
mv $REPO /go/src/$SITE/$ORG/$REPO
cd /go/src/$SITE/$ORG/$REPO
echo "Building"

go build -o /$OUT_DIR/$REPO
echo "Done building $SITE/$ORG/$REPO"
