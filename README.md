# pasty
Pasty is a fast and lightweight code pasting server

## General environment variables
| Environment Variable          | Default Value | Type     | Allowed Values  | Description                                                                                                 |
|-------------------------------|---------------|----------|-----------------|-------------------------------------------------------------------------------------------------------------|
| `PASTY_WEB_ADDRESS`           | `:8080`       | `string` | any             | Defines the address the webs erver listens to                                                               |
| `PASTY_STORAGE_TYPE`          | `file`        | `string` | `file`          | Defines the storage type the pastes are saved to                                                            |
| `PASTY_HASTEBIN_SUPPORT`      | `false`       | `bool`   | `true`, `false` | Defines whether or not the `POST /documents` endpoint should be enabled, as known from the hastebin servers |
| `PASTY_DELETION_TOKEN_LENGTH` | `12`          | `number` | any             | Defines the length of the deletion token of a paste                                                         |
| `PASTY_RATE_LIMIT`            | `30-M`        | `string` | any             | Defines the rate limit of the API (see https://github.com/ulule/limiter#usage)                              |

## Storage types
Pasty supports multiple storage types, defined using the `PASTY_STORAGE_TYPE` environment variable.
Every single one of them has its own configuration variables: 

### File
| Environment Variable      | Default Value | Type     | Allowed Values | Description                                               |
|---------------------------|---------------|----------|----------------|-----------------------------------------------------------|
| `PASTY_STORAGE_FILE_PATH` | `./data`      | `string` | any            | Defines the file path the paste files are being stored to |