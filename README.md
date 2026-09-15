# Facebook Backend V2

Go microservices for a Facebook-like platform. Kong sits in front of the services.

Frontend: [facebook-fe-v2](https://github.com/BT2701/facebook-fe-v2)

## Version

**0.1.0** — shared bootstrap, health checks, graceful shutdown, Docker Compose that actually waits for Mongo/Redis.

## Services

| Service | Port | Responsibility |
| --- | --- | --- |
| user-service | 8080 | Auth, accounts, JWT |
| notification-service | 8081 | Activity notifications |
| chat-service | 8082 | Messages and WebSocket |
| media-service | 8083 | Image upload and static files |
| post-service | 8084 | Posts, stories, comments, reactions |
| friend-service | 8085 | Friend graph and requests |
| game-service | 8086 | Slot game sessions, results, player balance |
| Kong | 8000 | API gateway |

## Stack

Go 1.23 · Echo · MongoDB · Redis · JWT · Kong · Docker Compose

Shared code lives in `shared/` (`config`, `httpx`). Each service keeps its own domain and adapters.

## Run locally

Prerequisites: Docker, Go 1.23+.

```sh
cp .env.example .env
# copy or edit each service .env if you run binaries outside Docker

docker compose up --build
```

Gateway: `http://localhost:8000`

Health:

```sh
make health
# or
curl http://localhost:8080/health
```

Frontend should use `REACT_APP_API_URL=http://localhost:8000`.

## Layout

```
shared/                 # env load, CORS, /health, graceful shutdown
user-service/
post-service/
friend-service/
chat-service/
notification-service/
media-service/
game-service/
api-gateway/kong.yml
```

Each service follows inbound adapter → application service → outbound repository.

## Notes

- SMTP for forgot-password is read from `SMTP_FROM` and `SMTP_PASSWORD`. Do not hardcode mail credentials.
- Destructive wipe routes (`DELETE /users`, `DELETE /posts`, …) are no longer public.
- Kong runs DB-less from `api-gateway/kong.yml`. There is no migration container.

## License

Apache License 2.0. See [LICENSE](LICENSE).
