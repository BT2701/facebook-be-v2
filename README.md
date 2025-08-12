# Facebook Backend V2 - Golang

This project is a backend service for a Facebook-like application, built using modern technologies and a microservices architecture.

## Link to Front-end: [Click here](https://github.com/BT2701/facebook-fe-v2)

---

## Project Overview

### Version
- **Current Version**: 0.0.2

### Services
The project is divided into multiple microservices, each responsible for a specific domain within the application:

1. **User Service**  
   Manages users, accounts, authentication, and authorization (using JWT tokens).

2. **Post Service**  
   Manages posts, stories, interactions, comments, and reactions; integrated with Google Drive Service for media storage.

3. **Friend Service**  
   Manages friend relationships and handles friend requests.

4. **Chat Service**  
   Manages messaging and calls (audio, video) using WebSocket and SignalR.

5. **Notification Service**  
   Handles notifications related to friend requests, post interactions, and other activities (via WebSocket).

6. **Game Service**  
   Manages games available on the platform; currently supports slot games and plans to expand to casual, arcade, and fishing games to diversify gameplay.

7. **API Gateway**  
   Routes and manages communication between services using **Kong API Gateway**.

---

## Technologies Used

- **Golang**: Main programming language for backend services.
- **Echo**: High-performance, minimalist Go web framework.
- **JWT**: Secure authentication and authorization.
- **MongoDB**: NoSQL database for storing user and application data.
- **Redis**: Caching layer for improved performance and reduced database load.
- **Docker**: Containerization of all services for easy deployment.
- **Kong Gateway**: API Gateway for routing and managing microservices traffic.
- **REST APIs**: For service-to-service and client communication.
- **WebSocket**: Real-time communication for chat, notifications, and games.
- **SignalR**: Real-time signaling for audio and video calls.

---

## Version Control

- **Git/GitHub**: Version control and code hosting.

---

## Development Tools

- **Visual Studio Code**: Main editor for development.
- **Docker Compose**: For running the system locally with all dependencies.

---

## Getting Started

### Prerequisites

- Docker installed on your machine.
- Golang installed (v1.20+ recommended).
- MongoDB and Redis instances running (can also run via Docker).

### Running the Application

1. Clone the repository:
    ```sh
    git clone https://github.com/BT2701/facebook-be-v2.git
    ```
2. Navigate to the project directory:
    ```sh
    cd facebook-be-v2
    ```
3. Build and run using Docker:
    ```sh
    docker-compose up --build
    ```

---

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

---

## License

This project is licensed under [Apache License Version 2.0](LICENSE).
