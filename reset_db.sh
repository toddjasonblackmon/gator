#! /usr/bin/bash

pushd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" down
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
popd
