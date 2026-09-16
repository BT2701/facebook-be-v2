.PHONY: tidy up down logs health

tidy:
	go work sync
	cd shared && go mod tidy
	cd user-service && go mod tidy
	cd post-service && go mod tidy
	cd friend-service && go mod tidy
	cd chat-service && go mod tidy
	cd notification-service && go mod tidy
	cd media-service && go mod tidy
	cd game-service && go mod tidy
	cd ai-service && go mod tidy

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f --tail=100

health:
	@curl -sf http://localhost:8000/user/health && echo user-service
	@curl -sf http://localhost:8000/notification/health && echo notification-service
	@curl -sf http://localhost:8000/chat/health && echo chat-service
	@curl -sf http://localhost:8000/media/health && echo media-service
	@curl -sf http://localhost:8000/post/health && echo post-service
	@curl -sf http://localhost:8000/friend/health && echo friend-service
	@curl -sf http://localhost:8000/game/health && echo game-service
	@curl -sf http://localhost:8000/ai/health && echo ai-service
