# Spy Cat API

This is a simple CRUD application for managing spy cats, missions, and targets. It demonstrates building RESTful APIs, interacting with a SQL-like database, and integrating with third-party services.

## Table of Contents
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Database Schema](#database-schema)

## Features

- Create, update, list, and delete spy cats
- Create, update, list, and delete missions
- Create, update, list, and delete targets
- Assign cats to missions and manage targets within missions
- Integrate with TheCatAPI to validate cat breeds

## Requirements

- Go 1.20+
- Docker
- Docker Compose

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/TarkovskyiTaras/Spy-Cat
    cd Spy-Cat
    ```

2. Install Go dependencies:

    ```sh
    go mod download
    ```

3. Connect to the database via GUI to see uuids:

    ```sh
    postgres://postgres:123456@localhost:5777/spycat_test?sslmode=disable
    ```

## API Endpoints

```sh
    make start all
```

## API Endpoints

### Cats

- **Create Cat**: `POST /cat`
    - Request Body: 
        ```json
        {
            "name": "Whiskers",
            "years_of_experience": 3,
            "breed": "Abyssinian",
            "salary": 5000.00
        }
        ```
    - Response: `201 Created`

- **List Cats**: `GET /cats/all`
    - Response: `200 OK`
    ```json
    [
        {
            "id": "uuid",
            "name": "Whiskers",
            "years_of_experience": 3,
            "breed": "Abyssinian",
            "salary": 5000.00
        }
    ]
    ```

- **Get Cat**: `GET /cat/{id}`
    - Response: `200 OK`
    ```json
    {
        "id": "some_uuid",
        "name": "Whiskers",
        "years_of_experience": 3,
        "breed": "Abyssinian",
        "salary": 5000.00
    }
    ```

- **Update Cat**: `PUT /cat/{id}`
    - Request Body:
        ```json
        {
            "salary": 6000.00
        }
        ```
    - Response: `200 OK`

- **Delete Cat**: `DELETE /cat/{id}`
    - Response: `200 OK`

### Missions

- **Create Mission**: `POST /mission`
    - Request Body:
        ```json
        {
            "name": "Mission Alpha",
            "cat_id": "some_uuid",
            "targets": [
                {
                    "name": "Target 1",
                    "country": "Country 1",
                    "notes": "Notes for Target 1"
                }
            ]
        }
        ```
    - Response: `201 Created`

- **List Missions**: `GET /missions`
    - Response: `200 OK`

- **Get Mission**: `GET /mission/{id}`
    - Response: `200 OK`

- **Update Mission**: `PUT /mission/{id}`
    - Request Body:
        ```json
        {
            "complete": true
        }
        ```
    - Response: `200 OK`

- **Delete Mission**: `DELETE /mission/{id}`
    - Response: `200 OK`

### Targets

- **Add Target**: `POST /target/{mission_id}`
    - Request Body:
        ```json
        {
            "name": "Target 2",
            "country": "Country 2",
            "notes": "Notes for Target 2"
        }
        ```
    - Response: `201 Created`

- **Get Target**: `GET /targets/{id}`
    - Response: `200 OK`

- **Update Target**: `PUT /targets/{id}`
    - Response: `200 OK`

- **Delete Target**: `DELETE /targets/{id}`
    - Response: `200 OK`

## Database Schema

```sql
CREATE TABLE cats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL,
    breed VARCHAR(255) NOT NULL,
    salary DECIMAL(10, 2) NOT NULL
);

CREATE TABLE missions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    cat_id UUID,
    complete BOOLEAN NOT NULL,
    FOREIGN KEY (cat_id) REFERENCES cats(id) ON DELETE SET NULL
);

CREATE TABLE targets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mission_id UUID NOT NULL,
    name VARCHAR(255),
    country VARCHAR(255),
    notes TEXT,
    complete BOOLEAN NOT NULL,
    FOREIGN KEY (mission_id) REFERENCES missions(id) ON DELETE CASCADE
);