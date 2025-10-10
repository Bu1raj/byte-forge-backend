# LOCAL SET UP

## Prerequisites
- [Go install](https://go.dev/)
- Docker desktop
## Setup instructions
Clone the repository and run the below commad from root of repository:
```
go mod tidy
```
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







