# API

> **IMPORTANT:** Even though the API is defined pretty well, it may encounter breaking changes while on the `develop` branch!

The REST API provided by pasty is the most important entrypoint when it comes to interacting with it. Basically everything, including the pasty frontend, is built on top of it.
To make things easier for other developers who decide to develop something in connection to pasty, everything important about it is documented here.

## Authentication/Authorization

Not everyone should be able to view, edit or delete all pastes. However, admins should be.
In order to achieve that, an effective auth flow is required.

There are two ways of authenticating:

### 1.) Paste-pecific

The `Authorization` header is set to `Bearer <modification_token>`, where `<modification_token>` is replaced with the corresponding paste-specific **modification token**.
This authentication is only valid for the requested paste.

### 2.) Admin tokens

The `Authorization` header is set to `Bearer <admin_token>`, where `<admin_token>` is replaced with the configured **administration token**.
This authentication is valid for all endpoints, regardless of the requested paste.

### Notation

In the folllowing, all endpoints that require an **admin token** are annotated with `[ADMIN]`.
All endpoints that are accessible through the **admin and modification token** are annotated with `[PASTE_SPECIFIC]`.
All endpoints that are accessible to everyone are annotated with `[UNSECURED]`.

## The paste entity

The central paste entity has the following fields:

* `id` (string)
* `content` (string)
* `modificationToken` (string)
    * The token used to authenticate with paste-specific secured endpoints; stored hashed and only returned on initial paste creation
* `created` (int64; UNIX timestamp)
* `metadata` (key-value store)
    * Different frontends may store simple key-value metadata pairs on pastes to enable specific functionality (for example clientside encryption)

### Encryption

The frontend pasty ships with implements an encryption option. This en- and decrypts pastes clientside and appends the HEX-encoded en-/decryption key to the paste URL (after a `#` because the so called **hash** is not sent to the server).
If a paste is encrypted using this feature, its `metadata` field contains a field like this:

```jsonc
{
    // --- omitted other entity field
    "metadata": {
        "pf_encryption": {
            "alg": "AES-CBC",                           // The algorithm used to encrypt the paste (currently, only AES-CBC is used)
            "iv": "54baa80cd8d8328dc4630f9316130f49"    // The HEX-encoded initialization vector of the AES-CBC encryption
        }
    }
}
```

## Endpoints

### [UNSECURED] Retrieve application information

```http
GET /api/v2/info
```

**Request:**
none

**Response:**
```json
{
    "modificationTokens": true,
    "reports": true,
    "version": "dev"
}
```

---

### [UNSECURED] Retrieve a paste

```http
GET /api/v2/pastes/{paste_id}
```

**Request:**
none

**Response:**
```json
{
    "id": "paste_id",
    "content": "paste_content",
    "created": 0000000000,
    "metadata": {},
}
```

---

### [UNSECURED] Create a paste

```http
POST /api/v2/pastes
```

**Request:**
```jsonc
{
    "content": "paste_content", // Required
    "metadata": {}              // Optional
}
```

**Response:**
```json
{
    "id": "paste_id",
    "content": "paste_content",
    "modificationToken": "raw_modification_token",
    "created": 0000000000,
    "metadata": {},
}
```

---

### [PASTE_SPECIFIC] Update a paste

```http
PATCH /api/v2/pastes/{paste_id}
```

**Request:**
```jsonc
{
    "content": "new_paste_content", // Optional
    "metadata": {}                  // Optional
}
```

**Response:**
```json
{
    "id": "paste_id",
    "content": "new_paste_content",
    "created": 0000000000,
    "metadata": {},
}
```

**Notes:**
* Changes in the `metadata` field only affect the corresponding field and don't override the whole key-value store (`{"metadata": {"foo": "bar"}}` will effectively add or replace the `foo` key but won't affect other keys).
* To remove a key from the key-value store simply set it to `null`.

---

### [PASTE_SPECIFIC] Delete a paste

```http
DELETE /api/v2/pastes/{paste_id}
```

**Request:**
none

**Response:**
none

---

### [UNSECURED] Report a paste

```http
POST /api/v2/pastes/{paste_id}/report
```

**Request:**
```json
{
    "reason": "reason"
}
```

**Response:**
```jsonc
{
    "success": true,        // Whether or not the report was received successfully (this is returned by the report webhook to allow custom errors)
    "message": "message"    // An optional message to display to the reporting user
}
```

**Notes:**
* The endpoint is only available if the report system is enabled. Otherwise it will return a `404 Not Found` error.
* The request that will reach the report webhook looks like this:
    ```json
    {
        "paste": "paste_id",
        "reason": "reason"
    }
    ```
* The response from this endpoint is the exact same that pasty expects from the webhook.