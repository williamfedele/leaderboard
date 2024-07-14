<div>
    <h3 align="center">Leaderboard</h3>
    <p align="center">
    A leaderboard using Redis Zset for educational purposes.
    </p>
</div>

## About

This was a little project used to learn about a few different tools. It uses Go as client/server, Redis for in-memory storage using its Zset, and gRPC for efficient client-server communication using protobuf for serialization.

## Getting Started

### Prerequisites

gRPC and the protobuf compiler need to be installed. The prerequisites can be followed [here](https://grpc.io/docs/languages/go/quickstart/).

## Usage

### Server

The protobuf compiler and the server can be started together using the `run` Make target.
```
> make run
Protobuf files compiled.
2024/07/14 15:56:23 Server started on port :8080
```

These can be run separately if needed using the `make compile_proto` and `make run_server` Make targets.

### Client

The demo client can also be run using Make target. The client will generate 100 random™ player scores in [0,100), print the top 5 scores, and print the players within a radius of 2 from Player7.
```
> make run_client
2024/07/14 15:57:55 Adding 100 random scores...
2024/07/14 15:57:55 Top Scores:
2024/07/14 15:57:55 1. Player: Player87, Score: 98
2024/07/14 15:57:55 2. Player: Player86, Score: 94
2024/07/14 15:57:55 3. Player: Player60, Score: 94
2024/07/14 15:57:55 4. Player: Player94, Score: 93
2024/07/14 15:57:55 5. Player: Player62, Score: 93
2024/07/14 15:57:55 Scores around Player7:
2024/07/14 15:57:55 1. Player: Player14, Score: 19
2024/07/14 15:57:55 2. Player: Player49, Score: 17
2024/07/14 15:57:55 3. Player: Player7, Score: 16
2024/07/14 15:57:55 4. Player: Player53, Score: 16
2024/07/14 15:57:55 5. Player: Player30, Score: 15
```

### Clean

There is a Make target to clean the output directory.
```
> make clean
Cleaned up generated files.
```

## License

Distributed under the MIT License. See `LICENSE` for more information.
