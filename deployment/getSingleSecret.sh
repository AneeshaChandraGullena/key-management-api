#!/bin/bash

# set -x

# cf target info below
# API endpoint:   https://api.stage1.ng.bluemix.net (API version: 2.54.0)
# User:           keymaster@bg.vnet.ibm.com
# Org:            keymaster@bg.vnet.ibm.com
# Space:          dev


TOKEN="bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI5YzI0MWU0OC0xMGU2LTQzMzUtYTZmZi0yNzM5ZTQxOTc1YTYiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJwYXNzd29yZC53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiLCJvcGVuaWQiLCJ1YWEudXNlciJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiYzU0N2ZmMDEtYThiOC00Mzg3LWJlNTYtYjRiMjYyYmJkYzE0Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImVtYWlsIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsInJldl9zaWciOiIyOGVjM2I3NyIsImlhdCI6MTQ4ODMxNDQ4NSwiZXhwIjoxNDg5NTI0MDg1LCJpc3MiOiJodHRwczovL3VhYS5zdGFnZTEubmcuYmx1ZW1peC5uZXQvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsiY2xvdWRfY29udHJvbGxlciIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiXX0.vAs-jFhfi1tIWwCndJDtWipXhSZnOgvgBHnrJ5RRZ3I"
SPACE="c6c5c1fd-d63c-44a9-8aee-0eaee9723ff6"
KP_API_SERVER_IP="localhost"
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
}' "${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets" | python -m json.tool | grep "id" | cut -d ":" -f2 | sed 's/.$//' | awk '{ print $1 }' | tr -d "\"")

echo $SECRET_ID
curl ${CURL_FLAGS} -X GET -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" \
"${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets/${SECRET_ID}" | python -m json.tool
