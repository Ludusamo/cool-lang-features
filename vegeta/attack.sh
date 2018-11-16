#!/bin/bash

if [ ! -d tmp ]; then
    mkdir -p tmp
fi

go install ../backend/backend.go
backend --listen 9999&
BACKGROUND_PID=$!
go install ../frontend/frontend.go
frontend --listen 9998 --backend localhost:9999&
FOREGROUND_PID=$!
python3 attack_gen.py
vegeta attack -rate=1000 -duration=30s -targets=tmp/vegeta_attack | vegeta report

echo ${BACKGROUND_PID}
echo ${FOREGROUND_PID}
rm -rf tmp
kill ${BACKGROUND_PID}
kill ${FOREGROUND_PID}
