# Central API

## Usage
Launch the API the following way:
```
go build -o main src/*.go
./main
```
or
```
go run src/*.go
```

# Endpoints
This is where we describe the API endpoints and how they react to certain data.


## bump

**URL** Structure:
```
http://localhost:8080/V1/bump
```

Method: **POST**

**JSON** Body:
```json
{"guildId": "656234546654"}
```

**Responses**:

Code | Response Header | JSON Response | Info
--- | --- | --- | ---
**200** | `200 - Ok` | `{"guildId":65654,"timestamp":16504}` | **guildId** successfully bumped!
**200** | `200 - Added` | `{"guildId":65654,"timestamp":16504}` |  **guildId** added to database and successfully bumped!
**425** | `425 - TooEarly` | `{"guildId":65654,"timestamp":16504}` |  **guildId** bump delta hasn't exceeded the bumping interval. Try again later.
**400** | `400 - Bad Request` | *None* | Bad request, make sure there are no strings in **guildId**
**500** | `500 - InternalServerError` | *None* | Internal Server Error, try again later!

*Additional note*:
The **JSON Response** is always a direct and latest representation of the stored guild in the database. That means when getting a `200` or `425` status code, the `timestamp` attribute represents the time that guild was last bumped in a UNIX timestamp.
