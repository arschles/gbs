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

if [ -e "glide.yaml" ]; then
  echo "Fetching Glide dependencies"
  glide install
  if [ "$?" != "0" ]; then
    echo "Glide install failed"
    exit $?
  fi
fi

echo "Building"
if [ "$CROSS_COMPILE" == "1" ]; then
  gox -output="/$BIN_DIR/gbs_cross/{{.Dir}}_{{.OS}}_{{.Arch}}"
else
  go build -o $BIN_NAME
  mv ./$BIN_NAME $BIN_DIR/$BIN_NAME
fi

echo "Done building $SITE/$ORG/$REPO (moved to $BIN_DIR)"
