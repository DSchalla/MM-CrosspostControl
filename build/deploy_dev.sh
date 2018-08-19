#!/bin/bash
set -e
set -x
./package.sh;
rm ../../../mattermost/mattermost-server/plugins/com.dschalla.crosspostcontrol/crosspostcontrol;
cp crosspostcontrol ../../../mattermost/mattermost-server/plugins/com.dschalla.crosspostcontrol/crosspostcontrol;