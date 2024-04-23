#!/bin/sh
./api/kinko &

cd frontend
node build/index.js

wait
