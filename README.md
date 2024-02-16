# Wuphf.com
--------------

[![I just sent myself a Wuphf (The Office)](https://img.youtube.com/vi/8wfG8ngFvPk/0.jpg)](https://www.youtube.com/watch?v=8wfG8ngFvPk)


## Wuphf.com: Reborn in Go!

Welcome to the new and improved Wuphf.com, a digital world made possible by the magic of code (not to be confused with magic beans). This project is a modern reimplementation of the iconic website from The Office, built with Go microservices and orchestrated by Kubernetes.

### Architecture Overview:

Wuphf.com 2.0 is a microservice-based architecture leveraging the power of:

- **User service:** Manages user accounts, authentication, and profiles.
- **Notification service:** Handles sending instant notifications via Kafka.
- **API Gateway:** Acts as the single entry point, handles user authentication, routes requests to microservices, and communicates with Kafka.

Each service runs independently within its own container and communicates through gRPC and HTTP protocols. This modular approach ensures scalability, maintainability, and resilience.

### Progress
- [X] User Service
- [X] Notification Service
- [X] API Gateway
- [X] Docker Integration
- [X] K8s Integration
- [ ] Helm Chart
- [ ] Istio
- [ ] Tests
- [ ] CI
- [ ] Web UI
- [ ] React Native UI

### Technologies Used:

- **Go:** The backend language, known for its speed, efficiency, and growing popularity.
- **gRPC:** Enables efficient and secure inter-service communication.
- **HTTP:** Used for user-facing APIs and communication with the API Gateway.
- **Kafka:** A distributed streaming platform for reliable and scalable message delivery (the notification magic!).
- **Docker:** Packages microservices as self-contained units for easy deployment and scaling.
- **Docker Compose:** Simplifies multi-container application development and testing.
- **Kubernetes:** Manages containers in a production environment, automating deployment and scaling.

*Could have used RabbitMQ as well, but wanted to use Kafka*

### Get Involved:

If you're curious to explore the code, contribute, or simply relive the glory days of Wuphf.com, you're in the right place!

- **Clone the repository:** `git clone https://github.com/your-username/wuphf-dot-com-go.git`
- **Run locally:** Use docker-compose to start the development environment (`docker-compose up`).
- **Deploy to Kubernetes:** Follow the provided instructions to deploy the application to your Kubernetes cluster.
- **Contribute:** We welcome bug reports, feature requests, and pull requests!

### Remember:

- This is a fun side project, not an official product.
- Please use responsibly and don't annoy your coworkers with too many Woofs!

### Disclaimer:

- May contain traces of Scrantonicity.
- Not affiliated with Dunder Mifflin Paper Company or any real-world entities.

We hope you enjoy using and contributing to Wuphf.com 2.0! Let's build something more magical than Ryan's website ever could have been.
