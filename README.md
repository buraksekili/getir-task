## Local Development

```shell script
$ git clone https://github.com/buraksekili/getir-task.git
$ cd getir-task
$ go build
$ ./getir-task
```

## Endpoints

- The default `<ENDPOINT_URL>` for local development is `http://localhost:3000`. It can be modified through [`.env`](https://github.com/buraksekili/getir-task/blob/ee6b6151fe3ab70be0bbd35dafc5bcb32ea7f63b/.env) file.
- Also, `<ENDPOINT_URL>` is available to public through `https://bs-getir-task.herokuapp.com/`

**Caveat:** Do not forget to specify [GETIR_TASK_MONGO_URI](https://github.com/buraksekili/getir-task/blob/ee6b6151fe3ab70be0bbd35dafc5bcb32ea7f63b/.env#L1) before running the application on your local setup. If you are using public URL, you do not need to specify `GETIR_TASK_MONGO_URI`.

### Fetch Data from MongoDB

```shell script
POST /mongo HTTP/1.1
Content-Type: application/json
Accept: application/json
```

Request body:
```json
{
    "startDate" : "2016-01-26",
    "endDate"   : "2018-02-02",
    "minCount"  : 2700,
    "maxCount"  : 3000
}
```

```shell script
curl -s -X POST "<ENDPOINT_URL>/mongo" -H 'Content-Type: application/json' -d '{ "startDate": "2016-01-26", "endDate": "2018-02-02", "minCount": 2700, "maxCount": 3000}'
```


### Create a data in In-Memory

```shell script
POST /in-memory HTTP/1.1
Content-Type: application/json
Accept: application/json
```

Request body:
```json
{
    "key"   : "string",
    "value" : "string"
}
```

```shell script
curl -s -X POST "<ENDPOINT_URL>/in-memory" -H 'Content-Type: application/json' -d '{"key": "active-tabs","value": "getir"}'
```

### Get a data from In-Memory

```shell script
GET /in-memory?key=string HTTP/1.1
Accept: application/json
```

```shell script
curl -s -X GET "<ENDPOINT_URL>/in-memory?key=active-tabs"
```





