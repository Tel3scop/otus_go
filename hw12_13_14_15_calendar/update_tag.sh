#!/bin/bash

GITHUB_SHA=$(git rev-parse HEAD)
echo "GITHUB_SHA: $GITHUB_SHA"
TAG_NAME=$(echo $GITHUB_SHA | head -c5)
echo "TAG_NAME: $TAG_NAME"
awk -v tag="$TAG_NAME" '/^TAG=/ {$0="TAG="tag} {print}' ./.env > .env.tmp && mv .env.tmp .env

echo "Updated TAG to $TAG_NAME in .env file"
