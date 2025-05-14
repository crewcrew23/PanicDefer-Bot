## About
PanicDefer_bot - this is a bot that monitors the status of your services <br>
and warns if there is something wrong with the service <br>

## Features and plans
### Features:
- ping service and save ping stats
- you can see statsabout last ping

Plans:
- notification if had problems with service
- send warn if avg response time higher than usual


## How to run
### need
- postgres
- go 1.24.2 or higher
- docker (maybe)

The bot uses rabbitMQ, you can run mq in docker or on your machine

## run 
### local
- rename exemple.yaml in /config to local.yaml and change fields of config
- change db url in make file in command migrate
- exec `make run migrate`
- exec command in order state (before start run rabbitMQ)
  ```
  run-server
  run-handler
  run-worker
  ```
### Docker