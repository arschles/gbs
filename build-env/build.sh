#!/bin/bash

CLONE_URL="https://$SITE/$ORG/$REPO.git"

echo "cloning $CLONE_URL"
git clone $CLONE_URL
if [[ $? != 0 ]];then
  echo "error cloning"
  exit 1
fi

mkdir -p /go/src/$SITE/$ORG
mv $REPO /go/src/$SITE/$ORG/$REPO
cd /go/src/$SITE/$ORG/$REPO
echo "building"

go build -o /$OUT_DIR/$REPO
echo "done building $SITE/$ORG/$REPO"
