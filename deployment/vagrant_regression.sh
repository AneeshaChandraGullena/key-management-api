#!/bin/bash

#set -x

# cf target info below
# API endpoint:   https://api.stage1.ng.bluemix.net (API version: 2.54.0)
# User:           keymaster@bg.vnet.ibm.com
# Org:            keymaster@bg.vnet.ibm.com
# Space:          dev


TOKEN="Bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI1NjEyN2NjMC1lZTkyLTQ2ZmUtODA2NS1hYzk3NTgzY2M1ZjYiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJwYXNzd29yZC53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiLCJvcGVuaWQiXSwiY2xpZW50X2lkIjoiY2YiLCJjaWQiOiJjZiIsImF6cCI6ImNmIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9pZCI6ImM1NDdmZjAxLWE4YjgtNDM4Ny1iZTU2LWI0YjI2MmJiZGMxNCIsIm9yaWdpbiI6InVhYSIsInVzZXJfbmFtZSI6ImtleW1hc3RlckBiZy52bmV0LmlibS5jb20iLCJlbWFpbCI6ImtleW1hc3RlckBiZy52bmV0LmlibS5jb20iLCJyZXZfc2lnIjoiMjhlYzNiNzciLCJpYXQiOjE0NzM2OTI0NjMsImV4cCI6MTQ3NDkwMjA2MywiaXNzIjoiaHR0cHM6Ly91YWEuc3RhZ2UxLm5nLmJsdWVtaXgubmV0L29hdXRoL3Rva2VuIiwiemlkIjoidWFhIiwiYXVkIjpbImNsb3VkX2NvbnRyb2xsZXIiLCJwYXNzd29yZCIsImNmIiwib3BlbmlkIl19.P-Nfyk-Uo1B2dQ4FdeWh5D5EhYY2Fnlx6zInaT7918A"
SPACE="c6c5c1fd-d63c-44a9-8aee-0eaee9723ff6"
LOCAL="127.0.0.1"
DEV="169.54.117.95"
PRESTAGE="169.54.117.96"
STAGE="198.23.117.188"
KP_API_SERVER_IP=${STAGE}
BLUEMIX_ORG="keymaster@bg.vnet.ibm.com"
CURL_FLAGS="-k"
SCHEME="https"

echo "-------HEALTH CHECK-------"
# test admin healthcheck
curl ${CURL_FLAGS} -X GET -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}"  "${SCHEME}://${KP_API_SERVER_IP}:8990/admin/v1/healthcheck" | python -m json.tool


echo "-------GET ALL-------"
# get all secrets
OFFSET=2
LIMIT=12
curl ${CURL_FLAGS} -X GET -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" "${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets?offset=${OFFSET}&limit=${LIMIT}" | python -m json.tool

echo "-------POST-------"
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
}' "${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets" | python -m json.tool | grep "id" | cut -d ":" -f2 | sed 's/.$//' | awk '{ print $1 }' | tr -d "\"")

echo "secret ID: ${SECRET_ID}"

echo "-------GET-------"
curl ${CURL_FLAGS} -X GET -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" \
"${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets/${SECRET_ID}" | python -m json.tool

echo "-------DELETE-------"
curl ${CURL_FLAGS} -X DELETE -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" -H "Prefer: return=representation" \
"${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets/${SECRET_ID}" | python -m json.tool
