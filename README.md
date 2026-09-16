# Facebook Backend V2

Go microservices for a Facebook-like platform. Kong sits in front of the services.

Frontend: [facebook-fe-v2](https://github.com/BT2701/facebook-fe-v2)

## Version

**0.3.0** — gateway-only public surface, JWT on writes, Redis notification events, and timeout/retry on friend → user calls.

## Services

| Service | Port | Responsibility |
| --- | --- | --- |
| user-service | 8080 | Auth, accounts, JWT, user search |
| notification-service | 8081 | Inbox, mark read, delete by action |
| chat-service | 8082 | Messages and WebSocket |
| media-service | 8083 | Image upload and static files |
| post-service | 8084 | Posts, stories, comments, reactions, search |
| friend-service | 8085 | Friends, requests, suggestions |
| game-service | 8086 | Slot config, bets, get-or-create player, server-side winnings |
| Kong | 8000 | API gateway |

## Stack

Go 1.23 · Echo · MongoDB · Redis · JWT · Kong · Docker Compose

Shared code lives in `shared/` (`config`, `httpx`, `auth`, `httpc`, `events`). Each service keeps its own domain and adapters.

## Run locally

Prerequisites: Docker, Go 1.23+.

```sh
cp .env.example .env
# copy or edit each service .env if you run binaries outside Docker

docker compose up --build
```

Gateway: `http://localhost:8000`. Service ports stay on the Docker network.

Health:

```sh
make health
# or
curl http://localhost:8000/user/health
```

Frontend should use `REACT_APP_API_URL=http://localhost:8000`.

## Layout

```
shared/                 # env, CORS, JWT writes, HTTP client, Redis events, /health
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
- Kong runs DB-less from `api-gateway/kong.yml`. Extra routes (`/Request`, `/comment`, `/reaction`) match the React client.
- Login responses no longer include the password hash. New game players are created with balance `1000`.
- Write endpoints require `Authorization: Bearer <jwt>` except login/register/forgot/reset.
- Friend request, accept, like, and comment publish Redis events; notification-service persists them.
- Friend suggestions call user-service with timeout, retry, and an `exclude` list instead of loading every user.
- Kong rate-limits the public gateway. Admin API binds to `127.0.0.1:8001`.

## License

Apache License 2.0. See [LICENSE](LICENSE).
