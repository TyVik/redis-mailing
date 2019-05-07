# stupid fun project (ಠ‿ಠ)

In the original, the task was: need to subscribe to redis, and send the received messages via websocket.

## Usage

While use is very simple))

```sh
$ docker-compose up --build -d
...
Successfully tagged redis-mailing:latest
Creating redis-mailing_redis_1 ... done
Creating redis-mailing_redis-mailing_1 ... done

$ go run test/tester.go
<current date and time> OK!
```
