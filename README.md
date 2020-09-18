# pasty
Pasty is a fast and lightweight code pasting server

## General environment variables
| Environment Variable          | Default Value | Type     | Description                                                                                                 |
|-------------------------------|---------------|----------|-------------------------------------------------------------------------------------------------------------|
| `PASTY_WEB_ADDRESS`           | `:8080`       | `string` | Defines the address the web server listens to                                                               |
| `PASTY_STORAGE_TYPE`          | `file`        | `string` | Defines the storage type the pastes are saved to                                                            |
| `PASTY_HASTEBIN_SUPPORT`      | `false`       | `bool`   | Defines whether or not the `POST /documents` endpoint should be enabled, as known from the hastebin servers |
| `PASTY_ID_LENGTH`             | `6`           | `number` | Defines the length of the ID of a paste                                                                     |
| `PASTY_DELETION_TOKEN_LENGTH` | `12`          | `number` | Defines the length of the deletion token of a paste                                                         |
| `PASTY_RATE_LIMIT`            | `30-M`        | `string` | Defines the rate limit of the API (see https://github.com/ulule/limiter#usage)                              |

## AutoDelete
Pasty provides an intuitive system to automatically delete pastes after a specific amount of time. You can configure it with the following variables:

## Storage types
Pasty supports multiple storage types, defined using the `PASTY_STORAGE_TYPE` environment variable (use the value behind the corresponding title in this README).
Every single one of them has its own configuration variables: 
| Environment Variable             | Default Value | Type     | Description                                                                    |
|----------------------------------|---------------|----------|--------------------------------------------------------------------------------|
| `PASTY_AUTODELETE`               | `false`       | `bool`   | Defines whether or not the AutoDelete system should be enabled                 |
| `PASTY_AUTODELETE_LIFETIME`      | `720h`        | `string` | Defines the duration a paste should live until it gets deleted                 |
| `PASTY_AUTODELETE_TASK_INTERVAL` | `5m`          | `string` | Defines the interval in which the AutoDelete task should clean up the database |

### File (`file`)
| Environment Variable      | Default Value | Type     | Description                                               |
|---------------------------|---------------|----------|-----------------------------------------------------------|
| `PASTY_STORAGE_FILE_PATH` | `./data`      | `string` | Defines the file path the paste files are being saved to  |

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

---

### MongoDB (`mongodb`)
| Environment Variable                      | Default Value                              | Type     | Description                                                     |
|-------------------------------------------|--------------------------------------------|----------|-----------------------------------------------------------------|
| `PASTY_STORAGE_MONGODB_CONNECTION_STRING` | `mongodb://pasty:pasty@example.host/pasty` | `string` | Defines the connection string to use for the MongoDB connection |
| `PASTY_STORAGE_MONGODB_DATABASE`          | `pasty`                                    | `string` | Defines the name of the database to use                         |
| `PASTY_STORAGE_MONGODB_COLLECTION`        | `pastes`                                   | `string` | Defines the name of the collection to use                       |

---

### SQL (`sql`)
| Environment Variable       | Default Value | Type     | Description                                                                         |
|----------------------------|---------------|----------|-------------------------------------------------------------------------------------|
| `PASTY_STORAGE_SQL_DRIVER` | `sqlite3`     | `string` | Defines the driver to use for the SQL connection (`sqlite3`, `postgres` or `mysql`) |
| `PASTY_STORAGE_SQL_DSN`    | `./db`        | `string` | Defines the DSN to use for the SQL connection                                       |
| `PASTY_STORAGE_SQL_TABLE`  | `pasty`       | `string` | Defines the table name to use for the SQL connection                                |