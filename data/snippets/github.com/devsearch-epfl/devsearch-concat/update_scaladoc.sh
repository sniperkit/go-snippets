#!/usr/bin/env bash

set -x

REPO="devsearch-concat"
DOC_FOLDER="api"

if [ "$TRAVIS_REPO_SLUG" == "devsearch-epfl/$REPO" ] && [ "$TRAVIS_JDK_VERSION" == "oraclejdk8" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ] && [ "$TRAVIS_BRANCH" == "master" ]; then
    
    # Setup travis info 
    git config --global user.email "travis@travis-ci.org"
    git config --global user.name "travis-ci"

    # Clone documentation repository 
    git clone --quiet --branch=gh-pages https://${GH_TOKEN}@github.com/devsearch-epfl/devsearch-doc
     
    # Clear and make new dir
    TARGET_FOLDER=devsearch-doc/$REPO/scaladoc/
    rm -rf $TARGET_FOLDER
    mkdir -p $TARGET_FOLDER
    
    # Copy scaladoc over
    cp -rv $DOC_FOLDER/* $TARGET_FOLDER
    
    # Step into repo
    cd devsearch-doc/
    
    # Mark for add and commit
    git add $REPO/scaladoc/*
    git commit -m "Travis #$TRAVIS_BUILD_NUMBER: Scaladoc for $REPO at commit $TRAVIS_COMMIT"
    git push origin gh-pages
    
    # Clean up documentaion repo
    cd ..
    rm -rf devsearch-doc/
else
    echo "Conditions not met to update scaladoc"
fi
