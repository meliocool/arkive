# ARKIVE - Decentralized Photo Sharing Application

A backend service built with Go that enables users to upload and share photos. The unique aspect of this project is that it leverages the Interplanetary File System (IPFS) for decentralized file storage, moving away from traditional cloud-based solutions.

This project showcases core backend development skills including:
* **RESTful API Design:** Building clear and well-structured API endpoints.
* **User Authentication:** Secure user registration and login using JWTs and 6-Digit Email Verification with ZohoMail.
* **Database Management:** Interacting with a PostgreSQL database to store user and photo metadata.
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

2.  Set up your environment variables in a `.env` file at the root of the project.
    ```bash
    POSTGRES_DB=arkive
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=dbPW
    POSTGRES_HOST=db
    POSTGRES_PORT=5432
    
    EMAIL_SMTP_USER=your_email
    EMAIL_SMTP_PASS=your_email_pw
    EMAIL_SMTP_PORT=587
    EMAIL_SMTP_HOST=your_email_provider

    JWT_SECRET=your_super_secret_key_here

    IPFS_API_KEY=your_pinata_api_key
    IPFS_API_SECRET=your_pinata_secret_api_key
    ```

3.  Build and run the application with Docker Compose:
    ```bash
    docker-compose up --build
    ```

---

## API Endpoints

| Endpoint                   | Method   | Description                                                   | Protected |
|----------------------------|----------|---------------------------------------------------------------| --------- |
| `/health`                  | `GET`    | A simple health check to ensure the server is running.        | No        |
| `/users/register`          | `POST`   | Registers a new user account.                                 | No        |
| `/users/verify`            | `POST`   | Verifies a user's account with a 6-digit code sent via email. | No        |
| `/users/login`             | `POST`   | Authenticates a user and returns a JWT.                       | No        |
| `/photos`                  | `POST`   | Uploads a photo to IPFS and saves its metadata.               | Yes       |
| `/photos`                  | `GET`    | Lists all photos uploaded by the authenticated user.          | Yes       |
| `/photos/:photoId`         | `DELETE` | Deletes a photo from IPFS and the database.                   | Yes       |
| `/photos/:photoId/profile` | `POST`   | Sets a photo as the authenticated user's profile picture.     | Yes       |
| `/public/photos`           | `GET`    | Lists all photos in the application for public viewing.       | No        |
| `/users/:userId`           | `GET`    | Returns a public profile and all photos for a specific user.  | No        |

---

## Architecture

* **API or Handler Layer:** Handles all incoming HTTP requests and routes them to the appropriate handlers.
* **Service Layer:** Contains the business logic, such as hashing passwords, generating JWTs, and interacting with external services.
* **Repository Layer:** Abstracts the database interactions (e.g., saving user data, fetching photo metadata).
* **IPFS Integration:** A dedicated service for communicating with the Pinata API to pin and retrieve files.

## Deployment
* **The API is live and deployed in a Droplet instance with DigitalOcean, however i dont have the money to buy a proper domain lol
