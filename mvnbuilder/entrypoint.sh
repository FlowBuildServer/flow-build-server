#!/bin/bash

set -x
set -euo pipefail

echo $GIT_REPO
echo $SOURCE_BRANCH
echo $TARGET_BRANCH


git clone $GIT_REPO
cd ${GIT_REPO##*/}
git checkout $TARGET_BRANCH
git branch -vv
git merge origin/$SOURCE_BRANCH

mvn clean install
