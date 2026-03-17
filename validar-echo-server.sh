#!/bin/bash

TEST_MESSAGE="Probando ando 123"

docker run -d -i --name server_echo_tester --rm --network "tp0_testing_net" ubuntu
docker exec server_echo_tester apt update >/dev/null
docker exec server_echo_tester apt install netcat-openbsd -y >/dev/null

RESPONSE=$(echo -n "$TEST_MESSAGE" | docker exec -i server_echo_tester nc server 12345 2>/dev/null)

docker stop server_echo_tester

if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
    exit 0
else
    echo "action: test_echo_server | result: fail"
    exit 1
fi

