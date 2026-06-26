#!/bin/sh

mkdir -p test
echo "QBIT_USER=admin" >> test/env
echo "QBIT_PASS=adminadmin" >> test/env

# This should exit cleanly, not doing anything
echo "===== TEST #1 ====="
radarr_eventtype=Test go run main.go --log "" --env "./test/env"
echo

# This should warn about the unusual state, but should complete fine
echo "Test" >> test/src
echo "===== TEST #2 ====="
radarr_eventtype=Download radarr_moviefile_sourcepath=./test/src radarr_moviefile_path=./test/dst radarr_download_id=0 go run main.go --log "" --env "./test/env"
rm test/src test/dst
echo

# This should warn about the unusual state, but should complete fine
echo "Test" >> test/dst
echo "===== TEST #3 ====="
radarr_eventtype=Download radarr_moviefile_sourcepath=./test/src radarr_moviefile_path=./test/dst radarr_download_id=0 go run main.go --log "" --env "./test/env"
rm test/src test/dst
echo

# This is the expected state and should exit cleanly
echo "Test" >> test/dst
echo "Test" >> test/src
echo "===== TEST #4 ====="
radarr_eventtype=Download radarr_moviefile_sourcepath=./test/src radarr_moviefile_path=./test/dst radarr_download_id=0 go run main.go --log "" --env "./test/env"
rm test/src test/dst
echo

# This should fail
echo "===== TEST #5 ====="
radarr_eventtype=Download radarr_moviefile_sourcepath=./test/src radarr_moviefile_path=./test/dst radarr_download_id=0 go run main.go --log "" --env "./test/env"
rm -rf test
