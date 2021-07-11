# pasty
Pasty is a fast and lightweight code pasting server

## !!! Important deprecation notices !!!

> This version of pasty uses a new field name for the so far called `deletionToken`: `modificationToken`.
> Instances using **PostgreSQL** are **not affected** as a corresponding SQL migration will run before the first startup.
> If you however use **another storage driver** you may have to **update the entries** by hand or using a simple query, depending on your driver as I don't plan to ship migrations for every single storage driver.
> It may be important to know that the **data migrator has been upgraded** too. This may serve as a **convenient workaround** (export data (field will be renamed) and import data with changed field names again).

> Additionally, I changed the three `DELETION_TOKEN*`environment variables to their corresponding `MODIFICATION_TOKEN*` ones:
> - `DELETION_TOKENS` -> `MODIFICATION_TOKENS`
> - `DELETION_TOKEN_MASTER` -> `MODIFICATION_TOKEN_MASTER`
> - `DELETION_TOKEN_LENGTH` -> `MODIFICATION_TOKEN_LENGTH`
>
> Again, **the old ones will still work** because I do not want to jumble your configurations. However, **please consider updating** them to stay future-proof ^^.


## Support

As pasty is an open source project on GitHub you can open an [issue](https://github.com/lus/pasty/issues) whenever you encounter a problem or feature request.
However, it may be annoying to open an issue just to ask a simple question about pastys functionalities, get help with the the installation process or mention something about the hosted version.
This is why I created a simple [Discord server](https://go.lus.pm/discord) you may want to join to get an answer to stuff like that pretty quickly.

## Disclaimer

The pasty web frontend comes with some service-related links in it (Discord server). Of course you are allowed to remove these references.
However, a small reference to pasty would be nice ^^.

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
| Environment Variable              | Default Value | Type     | Description                                                                                                                |
|-----------------------------------|---------------|----------|----------------------------------------------------------------------------------------------------------------------------|
| `PASTY_WEB_ADDRESS`               | `:8080`       | `string` | Defines the address the web server listens to                                                                              |
| `PASTY_STORAGE_TYPE`              | `file`        | `string` | Defines the storage type the pastes are saved to                                                                           |
| `PASTY_HASTEBIN_SUPPORT`          | `false`       | `bool`   | Defines whether or not the `POST /documents` endpoint should be enabled, as known from the hastebin servers                |
| `PASTY_ID_LENGTH`                 | `6`           | `number` | Defines the length of the ID of a paste                                                                                    |
| `PASTY_MODIFICATION_TOKENS`       | `true`        | `bool`   | Defines whether or not modification tokens should be generated                                                             |
| `PASTY_MODIFICATION_TOKEN_MASTER` | `<empty>`     | `string` | Defines the master modification token which is authorized to modify every paste (even if modification tokens are disabled) |
| `PASTY_MODIFICATION_TOKEN_LENGTH` | `12`          | `number` | Defines the length of the modification token of a paste                                                                    |
| `PASTY_RATE_LIMIT`                | `30-M`        | `string` | Defines the rate limit of the API (see https://github.com/ulule/limiter#usage)                                             |
| `PASTY_LENGTH_CAP`                | `50000`       | `number` | Defines the maximum amount of characters a paste is allowed to contain (a value `<= 0` means no limit)                     |

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