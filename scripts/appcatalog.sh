#!/bin/bash

set -x

./appcatalog $* &

wait_for_apis () {
    export API_BASE_URL=http://localhost:8081/api
    # wait for api
    ret=1
    while [ $ret -ne 0 ]
    do
        curl -s $API_BASE_URL/unsecured/version >/dev/null
        ret=$?
        echo "Return:" $ret
        sleep 1
    done
}

# Activate IMPORT_SAMPLE_DATA env var to inject sample data
[ "x${IMPORT_SAMPLE_DATA}" != "x" ] && {
    wait_for_apis
    ./mycompany.sh >/dev/null
}

# Any script called init-script.sh will be executed
# Use it for your own integration
[ -f ./init-script.sh ] && {
    wait_for_apis
    chmod 700 ./init-script.sh && ./init-script.sh
}

wait %1

