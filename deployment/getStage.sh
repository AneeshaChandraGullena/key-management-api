#!/bin/bash

# set -x

# cf target info below
# API endpoint:   https://api.stage1.ng.bluemix.net (API version: 2.54.0)
# User:           keymaster@bg.vnet.ibm.com
# Org:            keymaster@bg.vnet.ibm.com
# Space:          dev


TOKEN="Bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiIwOTlkZmNlYi1mM2JkLTRjZjUtOWRkYS01ODk0Y2Q2ZGFhNGEiLCJzdWIiOiI2ZjdhYmI0Ny1kZGJjLTRiZTgtYjNkMC1lMWIxOGEyZjhhOTAiLCJzY29wZSI6WyJvcGVuaWQiLCJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJwYXNzd29yZC53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiXSwiY2xpZW50X2lkIjoiY2YiLCJjaWQiOiJjZiIsImF6cCI6ImNmIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9pZCI6IjZmN2FiYjQ3LWRkYmMtNGJlOC1iM2QwLWUxYjE4YTJmOGE5MCIsIm9yaWdpbiI6InVhYSIsInVzZXJfbmFtZSI6InphYy5uaXhvbkBpYm0uY29tIiwiZW1haWwiOiJ6YWMubml4b25AaWJtLmNvbSIsImF1dGhfdGltZSI6MTQ3MzM3Mjc4MiwicmV2X3NpZyI6IjRlODZlY2IyIiwiaWF0IjoxNDczMzcyNzgyLCJleHAiOjE0NzQ1ODIzODIsImlzcyI6Imh0dHBzOi8vdWFhLnN0YWdlMS5uZy5ibHVlbWl4Lm5ldC9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjZiIsIm9wZW5pZCIsImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCJdfQ.UyOmI7YuvbGdUWEUCBTNYttN8cdOOqQZRK2iGyxQEwQ"
SPACE="263c02be-e9f5-478d-b8e7-09a4aa0fa0b9"
KP_API_SERVER_IP="ibm-key-protect.stage1.edge.bluemix.net"
#KP_API_SERVER_IP="169.54.117.96"
BLUEMIX_ORG="keymaster@bg.vnet.ibm.com"
CURL_FLAGS="-k"
SCHEME="https"

SECRET_ID=$(curl ${CURL_FLAGS} -X POST -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" -H "Prefer: return=representation" -d '{
    "metaData": {
        "collectiontype": "application/vnd.ibm.kms.secret+json",
        "collectionTotal": 1
    },
    "resources": [
        {
            "type": "application/vnd.ibm.kms.secret+json",
            "name": "testa secret",
            "algorithmType": "AES",
            "description": "a testing thing",
            "payload": "My super secret secret"
        }
    ]
}' "${SCHEME}://${KP_API_SERVER_IP}/api/v2/secrets" | python -m json.tool | grep "id" | cut -d ":" -f2 | sed 's/.$//' | awk '{ print $1 }' | tr -d "\"")

echo $SECRET_ID
curl ${CURL_FLAGS} -X GET -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" \
"${SCHEME}://${KP_API_SERVER_IP}/api/v2/secrets/${SECRET_ID}" | python -m json.tool
