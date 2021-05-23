# pasty
Pasty is a fast and lightweight code pasting server

## Installation

You may set up pasty in multiple different ways. However, I won't cover all of them as I want to keep this documentation as clean as possible.

### 1.) Building from source
To build pasty from source make sure you have [Go](https://go.dev) installed.

1. Clone the repository:
```sh
git clone https://github.com/lus/pasty
```

2. Switch directory:
```sh
cd pasty/
```

3. Run `go build`:
```sh
go build -o pasty ./cmd/pasty/main.go
```

To configure pasty, simply create a `.env` file in the same directory as the binary is placed in.

To run pasty, simply execute the binary.

### 2.) Docker (recommended)
To run pasty with Docker, you should have basic understanding of it.

An example `docker run` command may look like this:
```sh
docker run -d \
    -p 8080:8080 \
    --name pasty \
    -e PASTY_AUTODELETE="true" \
    ghcr.io/lus/pasty:latest
```

Pasty will be available at http://localhost:8080.

---

## General environment variables
| Environment Variable          | Default Value | Type     | Description                                                                                                        |
|-------------------------------|---------------|----------|--------------------------------------------------------------------------------------------------------------------|
| `PASTY_WEB_ADDRESS`           | `:8080`       | `string` | Defines the address the web server listens to                                                                      |
| `PASTY_STORAGE_TYPE`          | `file`        | `string` | Defines the storage type the pastes are saved to                                                                   |
| `PASTY_HASTEBIN_SUPPORT`      | `false`       | `bool`   | Defines whether or not the `POST /documents` endpoint should be enabled, as known from the hastebin servers        |
| `PASTY_ID_LENGTH`             | `6`           | `number` | Defines the length of the ID of a paste                                                                            |
| `PASTY_DELETION_TOKENS`       | `true`        | `bool`   | Defines whether or not deletion tokens should be generated                                                         |
| `PASTY_DELETION_TOKEN_MASTER` | `<empty>`     | `string` | Defines the master deletion token which is authorized to delete every paste (even if deletion tokens are disabled) |
| `PASTY_DELETION_TOKEN_LENGTH` | `12`          | `number` | Defines the length of the deletion token of a paste                                                                |
| `PASTY_RATE_LIMIT`            | `30-M`        | `string` | Defines the rate limit of the API (see https://github.com/ulule/limiter#usage)                                     |
| `PASTY_LENGTH_CAP`            | `50000`       | `number` | Defines the maximum amount of characters a paste is allowed to contain (a value `<= 0` means no limit)             |

## AutoDelete
Pasty provides an intuitive system to automatically delete pastes after a specific amount of time. You can configure it with the following variables:
| Environment Variable             | Default Value | Type     | Description                                                                    |
|----------------------------------|---------------|----------|--------------------------------------------------------------------------------|
| `PASTY_AUTODELETE`               | `false`       | `bool`   | Defines whether or not the AutoDelete system should be enabled                 |
| `PASTY_AUTODELETE_LIFETIME`      | `720h`        | `string` | Defines the duration a paste should live until it gets deleted                 |
| `PASTY_AUTODELETE_TASK_INTERVAL` | `5m`          | `string` | Defines the interval in which the AutoDelete task should clean up the database |

## Storage types
Pasty supports multiple storage types, defined using the `PASTY_STORAGE_TYPE` environment variable (use the value behind the corresponding title in this README).
Every single one of them has its own configuration variables: 

### File (`file`)
| Environment Variable      | Default Value | Type     | Description                                               |
|---------------------------|---------------|----------|-----------------------------------------------------------|
| `PASTY_STORAGE_FILE_PATH` | `./data`      | `string` | Defines the file path the paste files are being saved to  |

---

### PostgreSQL (`postgres`)
| Environment Variable         | Default Value                            | Type     | Description                                                                         |
|------------------------------|------------------------------------------|----------|-------------------------------------------------------------------------------------|
| `PASTY_STORAGE_POSTGRES_DSN` | `postgres://pasty:pasty@localhost/pasty` | `string` | Defines the PostgreSQL connection string (you might have to add `?sslmode=disable`) |

---

### MongoDB (`mongodb`)
| Environment Variable                      | Default Value                           | Type     | Description                                                     |
|-------------------------------------------|-----------------------------------------|----------|-----------------------------------------------------------------|
| `PASTY_STORAGE_MONGODB_CONNECTION_STRING` | `mongodb://pasty:pasty@localhost/pasty` | `string` | Defines the connection string to use for the MongoDB connection |
| `PASTY_STORAGE_MONGODB_DATABASE`          | `pasty`                                 | `string` | Defines the name of the database to use                         |
| `PASTY_STORAGE_MONGODB_COLLECTION`        | `pastes`                                | `string` | Defines the name of the collection to use                       |

---

### S3 (`s3`)
| Environment Variable                 | Default Value | Type     | Description                                                                               |
|--------------------------------------|---------------|----------|-------------------------------------------------------------------------------------------|
| `PASTY_STORAGE_S3_ENDPOINT`          | `<empty>`     | `string` | Defines the S3 endpoint to connect to                                                     |
| `PASTY_STORAGE_S3_ACCESS_KEY_ID`     | `<empty>`     | `string` | Defines the access key ID to use for the S3 storage                                       |
| `PASTY_STORAGE_S3_SECRET_ACCESS_KEY` | `<empty>`     | `string` | Defines the secret acces key to use for the S3 storage                                    |
| `PASTY_STORAGE_S3_SECRET_TOKEN`      | `<empty>`     | `string` | Defines the session token to use for the S3 storage (may be left empty in the most cases) |
| `PASTY_STORAGE_S3_SECURE`            | `true`        | `bool`   | Defines whether or not SSL should be used for the S3 connection                           |
| `PASTY_STORAGE_S3_REGION`            | `<empty>`     | `string` | Defines the region of the S3 storage                                                      |
| `PASTY_STORAGE_S3_BUCKET`            | `pasty`       | `string` | Defines the name of the S3 bucket (has to be created before setup)                        |