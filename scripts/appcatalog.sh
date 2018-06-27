#!/bin/bash

set -ex

./appcatalog $* &

[ "x${IMPORT_SAMPLE_DATA}" != "x" ] && {
    export API_BASE_URL=http://localhost:8081/api
    # wait for api
    ret=1
    while [ $ret -ne 0 ]
    do
        curl -s $API_BASE_URL/unsecured/version >/dev/null
        ret=$?
        echo $ret
        sleep 1
    done
    ./mycompany.sh 2>/dev/null >/dev/null
}

wait %1
