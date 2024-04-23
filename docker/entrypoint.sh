#!/bin/sh
/app/api/kinko &

node frontend/index.js

wait
