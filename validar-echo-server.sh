#!/bin/bash
VALIDATION_MESSAGE="Validacion echo server"

RESPONSE=$(docker run --rm --network tp0_testing_net alpine:latest sh -c "echo '$VALIDATION_MESSAGE' | nc server 12345")

if [ "$RESPONSE" = "$VALIDATION_MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi