#!/bin/bash

payload=./predict-payload.json

curl -d @${payload} -H "Content-Type: application/json" -v http://localhost:8080/invocations
