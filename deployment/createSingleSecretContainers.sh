#!/bin/sh

# set -x

# cf target info below
# API endpoint:   https://api.stage1.ng.bluemix.net (API version: 2.54.0)
# User:           keymaster@bg.vnet.ibm.com
# Org:            keymaster@bg.vnet.ibm.com
# Space:          dev


TOKEN="bearer eyJhbGciOiJIUzI1NiJ9.eyJqdGkiOiI2NTMzMmZkOS1jNWY0LTQwMTYtYjk4NC05YWI3ZDkxN2IwNTUiLCJzdWIiOiJjNTQ3ZmYwMS1hOGI4LTQzODctYmU1Ni1iNGIyNjJiYmRjMTQiLCJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJwYXNzd29yZC53cml0ZSIsImNsb3VkX2NvbnRyb2xsZXIud3JpdGUiLCJvcGVuaWQiLCJ1YWEudXNlciJdLCJjbGllbnRfaWQiOiJjZiIsImNpZCI6ImNmIiwiYXpwIjoiY2YiLCJncmFudF90eXBlIjoicGFzc3dvcmQiLCJ1c2VyX2lkIjoiYzU0N2ZmMDEtYThiOC00Mzg3LWJlNTYtYjRiMjYyYmJkYzE0Iiwib3JpZ2luIjoidWFhIiwidXNlcl9uYW1lIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsImVtYWlsIjoia2V5bWFzdGVyQGJnLnZuZXQuaWJtLmNvbSIsInJldl9zaWciOiIyOGVjM2I3NyIsImlhdCI6MTQ4OTQyODc1MCwiZXhwIjoxNDkwNjM4MzUwLCJpc3MiOiJodHRwczovL3VhYS5zdGFnZTEubmcuYmx1ZW1peC5uZXQvb2F1dGgvdG9rZW4iLCJ6aWQiOiJ1YWEiLCJhdWQiOlsiY2xvdWRfY29udHJvbGxlciIsInBhc3N3b3JkIiwiY2YiLCJ1YWEiLCJvcGVuaWQiXX0.zEbP5DIsDyCV36CW_doC_7T2IVGesvCO-ho1tGL_Jx8"
SPACE="c6c5c1fd-d63c-44a9-8aee-0eaee9723ff6"
#KP_API_SERVER_IP="10.140.107.67"
KP_API_SERVER_IP="127.0.0.1"
BLUEMIX_ORG="6a65e749-6f32-4c50-b584-0b5a2d412be4"
CURL_FLAGS="-k"
SCHEME="https"

# the following are only needed if Authentication is diabled
BLUEMIX_ROLE="Manager"
USER_ID="put your email here"

curl -w '%{http_code}\n' ${CURL_FLAGS} -X POST -H "Content-Type: application/json" -H "Authorization: ${TOKEN}" -H "Bluemix-Space: ${SPACE}" -H "Bluemix-Org: ${BLUEMIX_ORG}" -H "Prefer: return=representation" -H "Bluemix-User-Role: ${BLUEMIX_ROLE}" -H "User-Id: ${USER_ID}" -d '{
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
}' "${SCHEME}://${KP_API_SERVER_IP}:8990/api/v2/secrets"
