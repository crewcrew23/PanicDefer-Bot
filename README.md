# 🚨 PanicDefer Bot - Service Monitoring Solution

![GitHub](https://img.shields.io/badge/Go-1.24.2+-blue)
![GitHub](https://img.shields.io/badge/PostgreSQL-Supported-green)
![GitHub](https://img.shields.io/badge/RabbitMQ-Required-orange)

## 🌟 About

**PanicDefer_bot** is your reliable service monitoring assistant that constantly watches over your infrastructure and alerts you when something goes wrong. Never miss a service outage again!

## 🔥 Features

### ✅ Implemented
- 🏓 Continuous service pinging with statistics collection
- 📅 Save ping statistic 
- 🔔 Real-time notifications when services go down
- 📊 Historical ping data visualization
- ⚠️ Warning system for abnormal response times


### 🚧 Planned Features
- temporary no

## 🛠️ System Requirements

| Component       | Version           |
|-----------------|-------------------|
| Go              | 1.24.2 or higher |
| PostgreSQL      | 12+              |
| RabbitMQ        | 3.8+             |
| Docker (optional)| 20.10+          |

## 🚀 Installation & Running

### 📦 Prerequisites
1. Install [Go](https://go.dev/dl/)
2. Install [PostgreSQL](https://www.postgresql.org/download/)
3. Install [RabbitMQ](https://www.rabbitmq.com/download.html) or use Docker

### 🏃‍♂️ Local Setup
1. **Configure the application**:
## Local
you need register new telegram bot and enter you bot token into config <br>

```
cp config/example.yaml config/local.yaml
nano config/local.yaml  # Edit configuration

#you can run rabbitmq on docker or on your machine
docker run -d \
  --name rabbitmq \
  -p 5672:5672 \  #cliet port
  -p 15672:15672 \  # Web-interface
  -e RABBITMQ_DEFAULT_USER=<username> \
  -e RABBITMQ_DEFAULT_PASS=<secretPassword> \
  rabbitmq:3-management

#insert your db_url in Makefile in command migrate
make migrate
make run-server
make run-handler
make run-worker
```

## Docker
you need register new telegram bot and enter you bot token into config <br>
docker-compose provide default credentials for DB and rabbitMQ <br>
if you need you can change them <br>
base config named docker.yalm is contained in ./config/docker.yalm <br>
if you changed credo from docker-compose you need change credo in ./config/docker.yalm <br>

if you renamed the config file, you need to specify the new path in the .env file <br>

after ``` docker-compose up --build -d ```
