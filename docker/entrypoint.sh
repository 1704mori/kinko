#!/bin/sh
./api/kinko &

node frontend/index.js

wait
