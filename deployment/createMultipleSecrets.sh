#!/bin/bash

# set -x

# cf target info below
# API endpoint:   https://api.stage1.ng.bluemix.net (API version: 2.54.0)
# User:           keymaster@bg.vnet.ibm.com
# Org:            keymaster@bg.vnet.ibm.com
# Space:          dev

TOKEN="Bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI0Y2JhY2FiOC1hMmI5LTQ0NDktOTQwZC0yYjg3ZDI1YWU1MWYiLCJzdWIiOiIzNjFmOWE2MS1mZDA2LTQ4ZWUtOTY1NS0xMTU4NGE3NzNlMWMiLCJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJwYXNzd29yZC53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiLCJvcGVuaWQiXSwiY2xpZW50X2lkIjoiY2YiLCJjaWQiOiJjZiIsImF6cCI6ImNmIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9pZCI6IjM2MWY5YTYxLWZkMDYtNDhlZS05NjU1LTExNTg0YTc3M2UxYyIsIm9yaWdpbiI6InVhYSIsInVzZXJfbmFtZSI6Imttc3RhZ2VAdXMuaWJtLmNvbSIsImVtYWlsIjoia21zdGFnZUB1cy5pYm0uY29tIiwicmV2X3NpZyI6ImJmYWJmYjgxIiwiaWF0IjoxNDcyNjU1NTUyLCJleHAiOjE0NzM4NjUxNTIsImlzcyI6Imh0dHBzOi8vdWFhLnN0YWdlMS5uZy5ibHVlbWl4Lm5ldC9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjbG91ZF9jb250cm9sbGVyIiwicGFzc3dvcmQiLCJjZiIsIm9wZW5pZCJdfQ.QSBvYwfPSNLe4P-pBvxGQv0Xn49iIoKJKqiHFOYPgF0"
SPACE="c6c5c1fd-d63c-44a9-8aee-0eaee9723ff6"
KP_API_SERVER_IP="169.54.117.96"
BLUEMIX_ORG="keymaster@bg.vnet.ibm.com"
CURL_FLAGS="-k"
SCHEME="https"

curl ${CURL_FLAGS} -X POST -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" -H "Prefer: return=representation" -d '{
    "metaData": {
        "collectiontype": "application/vnd.ibm.kms.secret+json",
        "collectionTotal": 2
    },
    "resources": [
        {
            "type": "application/vnd.ibm.kms.secret+json",
            "name": "testa secret",
            "AlgorithmType": "AES",
            "description": "a testing thing",
            "payload": "My super secret secret"
        },
        {
            "type": "application/vnd.ibm.kms.secret+json",
            "name": "test B secret",
            "AlgorithmType": "AES",
            "description": "another testing thing",
            "payload": "My super secret thingy2"
        }
    ]
}' "${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets"
