#!/bin/bash


echo $GIT_REPO
echo $SOURCE_BRANCH
echo $TARGET_BRANCH


git clone $GIT_REPO
cd ${GIT_REPO##*/}
git checkout $TARGET_BRANCH
git merge $SOURCE_BRANCH

mvn clean install
