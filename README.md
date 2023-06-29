# swassistant-backend

Custom backend for [SWAssistant](https://github.com/CalebQ42/SWAssistant). Extension of [stupid-backend](https://github.com/CalebQ42/stupid-backend)

## APIs

### Profiles

Character, vehicles, and minion profiles.

> POST: /profile/upload?key={api_key}?type={character|vehicle|minion}

Upload a profile.

Request Body:

```json
{
    // profile data
}
```

Note: Only allows up to 5MB of data. If over 5MB returns 413.

Response:

```json
{
    "id": "profile ID",
    "expiration": 0 // Unix time (Seconds) of expiration
}
```

> GET: /profile/{profile id}?key={api_key}

Get an uploaded profile.

Response:

```json
{
    "type": "character|vehicle|minion",
    // profile data minus uid
}
```

### Rooms

> GET: /rooms/list?key={api_key}&token={jwt_token}

Get a list of rooms your currently a part of.

Response:

```json
[
    {
        "id": "room ID",
        "name": "room name",
        "owner": "username"
    }
]
```

> POST: /rooms/new?key={api_key}&token={jwt_token}?name={room name}

Create a new room

Response:

```json
{
    "id": "room ID",
    "name": "room name"
}
```

> GET: /rooms/{room id}?key={api_key}&token={jwt_token}

Get info about a room.

```json
{
    "id": "room ID",
    "name": "room name",
    "owner": "username",
    "users": [
        "username"
    ],
    "profiles": [
        "profile uuids"
    ]
}
```
