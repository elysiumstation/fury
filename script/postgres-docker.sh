#!/bin/bash

# It is important that snapshotscopy{to|from} path is accessible at the same location in
# the container and outside of it. If you're using a custom fury home, you must call this script
# with FURY_HOME set to your custom fury home when starting the database.

if [ -n "$FURY_HOME" ]; then
        FURY_STATE=${FURY_HOME}/state
else
        if [[ "$OSTYPE" == "linux-gnu"* ]]; then
                FURY_STATE=~/.local/state/fury
        elif [[ "$OSTYPE" == "darwin"* ]]; then
                FURY_STATE="${HOME}/Library/Application Support/fury"
        else
                 echo "$OSTYPE" not supported
        fi
fi

SNAPSHOTS_COPY_TO_PATH=${FURY_STATE}/data-node/networkhistory/snapshotscopyto
SNAPSHOTS_COPY_FROM_PATH=${FURY_STATE}/data-node/networkhistory/snapshotscopyfrom

mkdir -p "$SNAPSHOTS_COPY_TO_PATH"
chmod 777 "$SNAPSHOTS_COPY_TO_PATH"

mkdir -p "$SNAPSHOTS_COPY_FROM_PATH"
chmod 777 "$SNAPSHOTS_COPY_FROM_PATH"

docker run --rm \
           -e POSTGRES_USER=fury \
           -e POSTGRES_PASSWORD=fury \
           -e POSTGRES_DB=fury \
           -p 5432:5432 \
           -v "$SNAPSHOTS_COPY_TO_PATH":"$SNAPSHOTS_COPY_TO_PATH":z \
           -v "$SNAPSHOTS_COPY_FROM_PATH":"$SNAPSHOTS_COPY_FROM_PATH":z \
           timescale/timescaledb:2.8.0-pg14
