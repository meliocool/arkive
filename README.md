# ARKIVE - Decentralized Photo Sharing Application

A backend service built with Go that enables users to upload and share allPhotos. The unique aspect of this project is that it leverages the InterPlanetary File System (IPFS) for decentralized file storage, moving away from traditional cloud-based solutions.

This project showcases core backend development skills including:
* **RESTful API Design:** Building clear and well-structured API endpoints.
* **User Authentication:** Secure user registration and login using JWTs and bcrypt.
* **Database Management:** Interacting with a PostgreSQL database to store user and photo metadata.
* **Concurrency:** Utilizing Go's goroutines and channels for efficient, concurrent operations.
* **Modern DevOps Practices:** Containerizing the application and database with Docker and orchestrating them with Docker Compose.

---

## Getting Started

### Prerequisites

* Go (v1.20 or newer)
* Docker & Docker Compose
* A free IPFS pinning service account (e.g., Pinata)

### Installation

1.  Clone this repository:
    ```bash
    git clone https://github.com/meliocool/arkive.git
    cd arkive
    ```

2.  Set up your environment variables.
    * `DB_CONNECTION_STRING`
    * `JWT_SECRET`
    * `PINATA_API_KEY`
    * `PINATA_SECRET_API_KEY`

3.  Build and run the application with Docker Compose:
    ```bash
    docker-compose up --build
    ```

---

## API Endpoints (Soon to be developed)

| Endpoint                  | Method | Description                                                                 |
| ------------------------- | ------ | --------------------------------------------------------------------------- |
| `/health`                 | `GET`  | A simple health check to ensure the server is running.                      |
| `/users/register`         | `POST` | Registers a new user with a username and password.                          |
| `/users/login`            | `POST` | Authenticates a user and returns a JWT.                                     |
| `/allPhotos`                 | `POST` | Uploads a photo to IPFS and saves its metadata to the database. (Protected) |
| `/allPhotos`                 | `GET`  | Lists all allPhotos uploaded by the authenticated user. (Protected)            |

---

## Architecture

* **API Layer:** Handles all incoming HTTP requests and routes them to the appropriate handlers.
* **Service Layer:** Contains the business logic, such as hashing passwords and generating JWTs.
* **Repository Layer:** Abstracts the database interactions (e.g., saving user data, fetching photo metadata).
* **IPFS Integration:** A dedicated service for communicating with the Pinata API to pin and retrieve files.
