# Coop backend

## Design

- This project is design to handle social media messages across different platforms. Messages are collected by a seperate service to make the system flexible to add new social media platform, but also to keep spesific logic for the given platform isolated. 
- The mastodon service requests messages from mastodon.social and posts them to the main API.
- The main API collects messages and stores them in a PostgreSQL database. The reason for this approach was to be able to filter for a given topic, and to limit the load on the mastodon api in case of enforced rate limit etc.
- The system uses syncronous communication among services over http. I'm most comfortable with http-protocol and was the main reason for choosing that. 
- The system is built using pulling against the Mastodon api. I was going to implement the ability to holding the http connection open and pushing new data to the client, but did not get time to do so. There is thus no ability to stream from the API. 
- The Test database are used by messages_test.go in internal/database in the API. It's only a part of the compose for simplicity. The test db is used to test the data access layer. 

## Stack
- Chi router
- PostgreSQL database
- goose for db migrations
- Models built with raw SQL and Go standard libary
- Makefile for automating formatting, static checking, testing, etc

<strong>Created by Torkel Aannestad</strong>

- [torkelaannestad.com](https://torkelaannestad.com)
- [Github](https://github.com/Torkel-Aannestad)

## Prerequsits
### Docker with compose available
- https://www.docker.com/get-started/

# Getting Started

## .env file
- set up a .env file with the four values in the env.sample

## Run Docker Compose
- The folloing command starts the projects and starts fetching requests from mastodon.
```shell
  Docker Compose up --build
```

## API
## healthcheck
- Both services has a "/v1/healthcheck" endpoint
- api on port 4000, mastodon on port 5000
```shell
  curl  http://localhost:4000/v1/healthcheck
```

### Request data from the api
- Gives the last 10 posts from mastodon. Gets new entries every 30 seconds. 
```shell
  curl  http://localhost:4000/v1/messages
```

### Post a new message
- copy following curl command and run it against the running service
```shell
BODY='{"external_id": "100", "author": "John McClane", "body": "Hans Gruber is a bad guy", "platform": "mastodon"}'
  curl -d "$BODY" http://localhost:4000/v1/messages
```

## Reset the database
- docker compose down
- delete the db volume

## Run code audit and testing
```shell
  make  audit/api
  make  audit/mastodon
```

# Improvements & limitations
- Input validation - This is an important feature for enforcing the contract between services and to give good error messages when experiencing errors. I wanted to create a shared validator packes between the services, but was out of time. 
- Testing. The current testing for the project is limited. 
- Add central error codes look up add to errorResponse
- Vendor dependecies. By vendoring dependecies we have a copy of the code and are not dependent on repos service dependecies to build. 
- update if message already exist
- graceful shutdown - Currently the api will not wait for inflight requests and will close the server regardless. 
- Handle HTTP 429 (Too Many Requests) from the mastodon api by adjusting wait time and increaced backoff.
- Remove code duplication and add core module for shared helper functions and logic. 

