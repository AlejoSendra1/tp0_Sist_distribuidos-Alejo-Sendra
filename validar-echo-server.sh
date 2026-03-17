#!/bin/bash

TEST_MESSAGE="Probando ando 123"

docker run -d -i --name server_echo_tester --rm --network "tp0_testing_net" alpine

RESPONSE=$(echo -n "$TEST_MESSAGE" | docker exec -i server_echo_tester nc server 12345 2>/dev/null)

docker stop server_echo_tester

if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi

exit 0
