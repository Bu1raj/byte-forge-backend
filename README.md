# LOCAL SET UP

## Prerequisites
- [Go install](https://go.dev/)
- Docker desktop

## Setup instructions
Clone the repository and run the below commad from root of repository:
```
go mod tidy
```

## Configuration

The application uses environment variables for configuration. You can set these in a `.env` file in the root directory or export them directly in your shell.

### Configuration Variables

Copy `.env.example` to `.env` and modify as needed:
```bash
cp .env.example .env
```

Available configuration options:
- `SERVER_PORT`: HTTP server port (default: 8080)
- `KAFKA_BROKER`: Kafka broker address (default: localhost:29092)
- `KAFKA_PRODUCER_TOPICS`: Comma-separated list of Kafka producer topics
- `KAFKA_CONSUMER_TOPICS`: Comma-separated list of Kafka consumer topics
- `REDIS_ADDRESS`: Redis server address (default: localhost:6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number (default: 0)

**Note:** The server and worker use different Kafka topic configurations:
- **Server**: Produces to `submissions`, consumes from `results`
- **Worker**: Produces to `results`, consumes from `submissions`

----
> [!NOTE]
> Start the docker desktop before running these commands

**Start Kafka and Zookeeper**
```
docker-compose up -d
```

**Exec into kafka container**
```
docker exec -it kafka bash
```

**Create a `submissions` and `results` topic**
```
kafka-topics --create --topic submissions --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```
```
kafka-topics --create --topic results --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

**List the topics for confirmation**
```
kafka-topics --list --bootstrap-server localhost:9092
```

**Start the sever in one terminal**
```
go run ./cmd/server/
```

**Start the worker in another terminalr**
```
go run ./cmd/worker/
```

**Test with the below curl request** (observe the worker terminal
```
curl -X POST http://localhost:8080/submit \
  -H "Content-Type: application/json" \
  -d '{
    "language": "python",
    "code": "print(\"Hello from inside Docker!\")",
    "timeout_seconds": 5
  }'
```

> [!NOTE]
> before doing the `curl` request run `docker pull python:3.11-alpine`, or increase the timeout to 45s


**Kafka UI is present on localhost:8081







