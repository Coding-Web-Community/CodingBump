# Central API

## Usage
Launch the API the following way:
```
cd src
go build -o main .
./main
```
or
```
cd src
go run .
```

# Endpoints
This is where we describe the API endpoints and how they react to certain data.


# bump

**URL** Structure:
```
http://localhost:8080/V1/bump
```

Method: **POST**

**JSON** Body:
```json
{"guildId": 636145886279237652}
```

## **Responses**:


### **200**
- *200 - Ok* | **guildId** successfully bumped!
```json
{
    "code": 200,
    "message": "Guild bumped",
    "payload": {
        "guildId": 636145886279237652,
        "timestamp": 1601832221
    }
}
```

- *200 - Added* | **guildId** successfully bumped!
```json
{
    "code": 200,
    "message": "Guild added and bumped",
    "payload": {
        "guildId": 636145886279237152,
        "timestamp": 1601832254
    }
}
```

### **400**
- *400 - BadRequest* | Request body contains invalid character, most likely a string or a non number
allowed characters: 0-9
```json
{
    "code": 400,
    "message": "Request body contains invalid character",
    "payload": {
        "guildId": 0,
        "timestamp": 0
    }
}
```
- *400 - BadRequest* | **guildId** needs to be 18 characters long
```json
{
    "code": 400,
    "message": "GuildId does not conform to 18 character long integer requirement",
    "payload": {
        "guildId": 636149237152,
        "timestamp": 0
    }
}
```

### **425**
- *425 - TooEarly* | **guildId** bump delta hasn't exceeded the bumping interval. Try again later
```json
{
    "code": 425,
    "message": "Guild bumped too early",
    "payload": {
        "guildId": 636145886279237152,
        "timestamp": 1601832360
    }
}
```

### **500**
- *500 - InternalServerError* | Unrecoverable internal server error. Try again later
```json
{
    "code": 425,
    "message": "",
    "payload": {
        "guildId": 0,
        "timestamp": 0
    }
}
```

*Additional note*:
The **payload** is always a direct and latest representation of the stored guild in the database. That means when getting a `200` or `425` status code, the `timestamp` attribute represents the time that guild was last bumped in a UNIX timestamp.

# fetch

**URL** Structure:
```
http://localhost:8080/V1/fetch
```

Method: **GET**

## **Responses**:


### **200**
- *200 - Ok* | **guildId** successfully bumped!
```json
{
    "code":200,
    "message":"Ok",
    "paypload":[
      {"guildId":636145886279237699,"timestamp":1602394289},
      {"guildId":636123886245557612,"timestamp":1602394230}
    ]
}
```

### **400**
- *400 - BadRequest*
```json
{
    "code":400,
    "message":"BadRequest",
    "paypload":[
      {}
    ]
  }
```
