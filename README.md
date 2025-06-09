# :pray: Prayer Buddies - Backend :pray:

CI: ![CI Badge](https://github.com/megarage9000/Prayer-Buddies/actions/workflows/ci.yaml/badge.svg)

Prayer Buddies - Backend is the backend system used for the Prayer Buddies application. This project handles user requests and server interactions with the Prayer Buddies database. Features include:

- Authentication using **JWT tokens** and password hashing using bcrypt
- ORM function using generated **SQLC** to communicate with the Prayer Buddies **PostgreSQL** Database
- **CRUD REST endpoints** that handle:
    - Registering and logging users in, responding with a JWT
    - Sending prayer requests, uploaded to Prayer Buddies PostgreSQL Database
    - Listing sent / received prayer requests from Prayer Buddies PostgreSQL Database

# Requirements

- Linux / UNIX operating system (This was run in WSL2)
- Go version 1.23.x
- Postgre SQL (This was using PostgreSQL 16.9)

# How To Run

# Features Implemented

# :hammer: TODO Checklist :hammer:

