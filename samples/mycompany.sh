#!/bin/bash

API_BASE_URL=${API_BASE_URL:-http:/localhost:8081/api}

set -ex
HOST=`echo $API_BASE_URL | awk -F/ '{print $3}'`
BASEPATH=`echo $API_BASE_URL | grep / | cut -d/ -f4-`
SCHEME=`echo $API_BASE_URL | awk -F: '{print $1}'`

curl --request DELETE \
  --url $API_BASE_URL/v1/applications/

curl --request DELETE \
  --url $API_BASE_URL/v1/deployments/

curl --request DELETE \
  --url $API_BASE_URL/v1/environments/

curl --request PUT \
  --url $API_BASE_URL/v1/environments/prodfr \
  --header 'content-type: application/json' \
  --data '{
	"name":"Prod France",
	"properties": {
		"region": "fr"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/environments/produk \
  --header 'content-type: application/json' \
  --data '{
	"name":"Prod United Kingdom",
	"properties": {
		"region": "gb"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/environments/prodde \
  --header 'content-type: application/json' \
  --data '{
	"name":"Prod Germany",
	"properties": {
		"region": "de"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/environments/staging \
  --header 'content-type: application/json' \
  --data '{
	"name":"Staging",
	"properties": {
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "human-resources",
	"name": "travel-api",
	"version": "latest",
	"tags": ["travel", "expense"],
	"manifest": {
		"authors": [
			{
				"name": "John Doe",
				"email": "john.doe@mycompany.com",
				"role": "MAINTAINER"
			},
			{
				"name": "Foobar",
				"email": "support@foobar.com",
				"role": "VENDOR"
			}
		],
		"repository": "https://git.mycompany.com/HR/travel.git",
		"description": "API to submit travel requests and declare expenses"
	},
	"properties": {
		"readme": "Some read me for version 1.0.0"
	}
}'

OPENAPI_SPEC='{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0",
    "title": "Travel API",
    "description": "API that allows to book travels",
    "termsOfService": "http://mycompany.com/travelapi/terms/",
    "contact": {
      "name": "Human Resources Team"
    },
    "license": {
      "name": "MIT"
    }
  },
	"host": "'"$HOST"'",
  "basePath": "'"/$BASEPATH"'",
  "schemes": [
    "'"$SCHEME"'"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/": {
      "get": {
        "description": "Returns the root of the API",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "A list of subroutes.",
            "schema": {
              "type": "object",
              "items": {
                "$ref": "#/definitions/Links"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Links": {
      "type": "object",
      "properties": {
        "_links": {
          "type": "array"
        }
      }
    }
  }
}'
curl --request POST \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/1.0.0/deploy/prodfr \
  --header 'content-type: application/json' \
  --data "{
		\"properties\":{
			\"openapi\": $OPENAPI_SPEC
		}
	}"

curl --request POST \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/1.0.0/deploy/produk \
  --header 'content-type: application/json' \
  --data "{}"


curl --request PUT \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/1.0.1 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "human-resources",
	"name": "travel-api",
	"version": "latest",
	"tags": ["travel", "expense"],
	"manifest": {
		"authors": [
			{
				"name": "John Doe",
				"email": "john.doe@mycompany.com",
				"role": "MAINTAINER"
			},
			{
				"name": "Foobar",
				"email": "support@foobar.com",
				"role": "VENDOR"
			}
		],
		"repository": "https://git.mycompany.com/HR/travel.git",
		"description": "API to submit travel requests and declare expenses"
	}
}'

curl --request POST \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/1.0.1/deploy/prodde \
  --header 'content-type: application/json' \
  --data "{}"

curl --request PUT \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/2.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "human-resources",
	"name": "travel-api",
	"version": "latest",
	"tags": ["travel", "expense"],
	"manifest": {
		"authors": [
			{
				"name": "John Doe",
				"email": "john.doe@mycompany.com",
				"role": "MAINTAINER"
			},
			{
				"name": "Foobar",
				"email": "support@foobar.com",
				"role": "VENDOR"
			}
		],
		"repository": "https://git.mycompany.com/HR/travel.git",
		"description": "API to submit travel requests and declare expenses"
	},
	"properties": {
		"readme": "# Some readme elements for version 2.0.0\n\nwith some **notes**"
	}
}'

curl --request POST \
  --url $API_BASE_URL/v1/applications/human-resources/travel-api/versions/2.0.0/deploy/staging \
  --header 'content-type: application/json' \
  --data "{
		\"properties\":{
			\"openapi\": $OPENAPI_SPEC
		}
	}"

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/api/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "billing",
	"name": "api",
	"version": "1.0.0",
	"manifest": {
		"description": "Billing API"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/archiver/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "billing",
	"name": "archiver",
	"version": "1.0.0",
	"manifest": {
		"description": "Archive old bills"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/currency-converter/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "billing",
	"name": "currency-converter",
	"version": "1.0.0",
	"manifest": {
		"description": "Convert currencies using public rates"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/antifraud/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "billing",
	"name": "antifraud",
	"version": "1.0.0",
  "tags": ["order"],
	"manifest": {
		"description": "Rules engine to advise fraud risk on a order"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/web/ecommerce-api/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "web",
	"name": "ecommerce-api",
	"version": "1.0.0",
	"manifest": {
		"description": "Ecommerce api"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/web/cms-api/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "web",
	"name": "cms-api",
	"version": "1.0.0",
	"manifest": {
		"description": "Content management api"
	}
}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/web/ui/versions/1.0.0 \
  --header 'content-type: application/json' \
  --data '{
	"domain": "web",
	"name": "cms-api",
	"version": "1.0.0",
	"manifest": {
		"description": "Public web application"
	}
}'


curl --request PUT \
  --url $API_BASE_URL/v1/badges/readme \
  --header 'content-type: application/json' \
  --data '{
	"title": "README.md",
	"type": "enum",
	"levels": 
	[
		{
			"id": "unset",
			"label": "Unknown", 
			"color": "lightgrey",
			"description": "No data about the availability of README.md",
			"isdefault": true
		},
		 
		{
			"id": "notrelevant",
			"label": "Not relevant", 
			"color": "grey",
			"description": "This application will not have a README.md"
		},
		{
			"id": "notfound",
			"label": "Not Found", 
			"color": "red",
			"description": "The README.md file was not found in the repository"
		},
		{
			"id": "tooshort",
			"label": "Too Short", 
			"color": "orange",
			"description": "The README.md has been found but it is too short to contain valuable information"
		},
		{
			"id": "exists",
			"label": "✔", 
			"color": "green",
			"description": "The README.md has been found in the git repository"
		}
	]}'

	curl --request PUT \
  --url $API_BASE_URL/v1/badges/license \
  --header 'content-type: application/json' \
  --data '{
	"title": "License",
	"type": "enum",
	"levels": 
	[
		{
			"id": "unset",
			"label": "Unknown", 
			"color": "lightgrey",
			"description": "No data about the availability of a LICENSE file",
			"isdefault": true
		},
		 
		{
			"id": "notrelevant",
			"label": "Not relevant", 
			"color": "grey",
			"description": "This application will not have a LICENSE file"
		},
		{
			"id": "notfound",
			"label": "Not Found", 
			"color": "red",
			"description": "The LICENSE file was not found in the repository"
		},
		{
			"id": "gpl",
			"label": "GPL", 
			"color": "orange",
			"description": "A GPL License has been detected in LICENSE file"
		},
		{
			"id": "bsd",
			"label": "BSD", 
			"color": "green",
			"description": "A BSD License has been detected in LICENSE file"
		}
	]}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/api/versions/1.0.0/badgeratings/readme \
  --header 'content-type: application/json' \
  --data '{
		"level": "exists",
		"comment": "340 lines"
	}'

curl --request PUT \
  --url $API_BASE_URL/v1/applications/billing/api/versions/1.0.0/badgeratings/license \
  --header 'content-type: application/json' \
  --data '{
		"level": "bsd",
		"comment": "confidence level 85%"
	}'
