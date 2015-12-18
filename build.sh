#!/bin/bash

echo "cloning"
git clone https://$SITE/$ORG/$REPO.git
mkdir -p /go/src/$SITE/$ORG
mv $REPO /go/src/$SITE/$ORG/$REPO
cd /go/src/$SITE/$ORG/$REPO
echo "building"
go build -o /pwd/$REPO
