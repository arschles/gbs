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

if [ "$CROSS_COMPILE" == "1" ]; then
  gox -output="/$OUT_DIR/$REPO/{{.Dir}}_{{.OS}}_{{.Arch}}"
else
  go build -o /$OUT_DIR/$REPO
fi

echo "Done building $SITE/$ORG/$REPO (moved to $BIN_DIR/$BIN_NAME)"
