# Backend manual
Go v 1.24 is used as a backend for the ObjectWaterfall project. SqLite is used as a database.

## Running
For the project running, Go 1.24 or later should be install in the system. 
The command "go run ." should be run in the root directory of the project for start the backend.  
When the backend has been run it receives requests on the 8888 port. 

## Terms
|Name    | Description |
|:------:| ------------|
| Worker | A unit which is spamming requests to an API endpoint |

## Endpoints
There are seven endpoints in the project:

- add
- start
- stop
- seed
- getWorkers
- getRunningWorkers
- logsWs

All endpoints (exept logsWs) return a response with status code and a json string of the form - 
if result is succeed
```
{
    "result" : "some object"
}
```
if result is error
```
{
    "error" : "Error message"
}
```

### /add

#### Definition
Receives a POST request for addition a new worker and its settings to the database.

Receives no parameters and a json body which contains the following fields:

|Name|data type|descriprion|
|:---|:---:|---|
| workerName | String | The name of a new worker |
|timer| Float | A duration of the requests sending in minutes |
|requestDellay| Float | A dellay between requests in seconds |
|random| Boolean | Should a data be taken randomly from the database for the sertain worker |
|writesNumberToSend| Integer | A number of records which should be taken from database for a single request |
|totalToSend| Integer | A max number of records for sending for the whole test duration |
| stopWhenTableEnds | Boolean | Should a test be stopped or start from the table's start when the data in the table is end (Makes no affection if "random" flag is true) |

##### Example 
URL 
```
http://localhost:8888/add
```
json body
```
{
	"workerName":"The Best Worker",
	"timer":10000,
	"requestDellay":1,
	"random":false,
	"writesNumberToSend":20,
	"totalToSend":100000,
	"StopWhenTableEnds":false
}
```
#### Response
The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    "result": "The worker has been added successfully"
}
```
Failed result returns 400 status code (bad request)
```
{
    "error": "error message"
}
```

-----------------------------------------------------------------------------------------

### /start

#### Definition
Receives a POST request for the starting a sertain worker.

Receives 1 following URL parameter:

|Name|descriprion|
|:---|---|
| id | A worker identificator |

Receives a json body which contains the following fields:

|Name|data type|descriprion|
|:---|:---:|---|
| host | String | A URL to be spammed |
| authModel | Object | An autherization model for getting token (new ways of auth will be added) |
| authModel.authUrl | String | A URL where auth happens |
| authModel.model | String | a json string which contains all what is needed for authorization on an API |

>**Warning - the authModel.mel is a string exactly. If you use Postman or other tools, you should wrap the whole json object with quates and use shielding In the backend it is parsed and treatment**

#### Example
URL 
```
http://localhost:8888/start?id=1
```
json body
```
{
    "host": "http://localhost:2305/path",
    "authModel": {
        "authUrl":"localhost:2305/login",
        "model":"{userNane:Usen, userPassword: p@ssw0rd}"
    }
}
```
#### Response

The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    // Will be fixed as {"result": obj/stirng}
    "workerId": workerId 
}
```
Failed result returns 400 status code (bad request)/500 status code (internal server error)
```
{
    "error": "error message"
}
```
or 409 status code (Conflict) if worker has been already started 

```
{
    "error": "The worker 'workern name' is running alredy"
}
```

-----------------------------------------------------------------------------------------

### /stop

#### Definition
Receives a GET request for the ptopping a sertain worker.

Receives 1 following URL parameter and no body:

|Name|descriprion|
|:---|---|
| id | A worker identificator in worker store (not the same as a worker id in database)|

#### Example

URL 
```
http://localhost:8888/stop?id=1
```
json body
```
no json body
```
#### Response

The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    "result": "Ok"
}
```
Failed result returns 400 status code (bad request)
```
{
    "error": "error message"
}
```

-----------------------------------------------------------------------------------------

### /seed

#### Definition
Receives a POST request for the making a test dummy data and storing it in the database.

Receives no parameters and a json body which contains the following fields:

|Name|data type|descriprion|
|:---|:---:|---|
| workerName | String | A worker name the dummy data will be making for |
| jStr | String | A json string which is a base for dummy data. Fields of this json are replaced with generated data |
| count | Integer | A count of records to be created for certain worker |

>**Warning - the jStr is a string exactly. If you use Postman or other tools, you should wrap the whole json object with quates and use shielding In the backend it is parsed and treatment**

#### Example
URL 
```
http://localhost:8888/seed
```
json body
```
{
    "workerName": "TestWorker",
    "jStr" :"{"result": 
    [
        {
            "message": "Hello, Tara! Your order number is: #67","phoneNumber": "901.802.4032",
            "phoneVariation": "+90 313 772 10 16",
            "status": "disabled",
            "name": 
            {
                "first": "Federico",
                "middle": "Kennedy",
                "last": "Armstrong"
                },
                "username": "Federico-Armstrong",
                "password": "aSzfK2V2XOgGO1H",
                "emails": [
                    "Emelie98@example.com","Jakayla1@gmail.com"],
                    "location": {
                        "street": "7213 Hagenes Cape",
                        "city": "Kettering",
                        "state": "North Carolina",
                        "country": "Cambodia",
                        "zip": "31336",
                        "coordinates": {
                            "latitude": 68.4004,
                            "longitude": -81.5312
                            }
                            },"website": "https://severe-assist.info/",
                            "domain": "quick-witted-default.org",
                            "job": {
                                "title": "Direct Paradigm Consultant",
                                "descriptor": "Global",
                                "area": "Infrastructure",
                                "type": "Engineer",
                                "company": "Hayes and Sons"},
                                "creditCard": {
                                    "number": "3529-0083-5264-1525",
                                    "cvv": 557,
                                    "issuer": "mastercard"},
                                    "uuid":"5b8c7a77-d3e1-4547-ab65-6236710bac18",
                                    "objectId": "6857c01cf688af1fc5c99132"}]}",
    "count" : 100000
}
```
#### Example for Postman
json body
```
{
    "workerName": "TestWorker",
    "jStr" :"{\"result\": [{\"message\": \"Hello, Tara! Your order number is: #67\",\"phoneNumber\": \"901.802.4032\",\"phoneVariation\": \"+90 313 772 10 16\",\"status\": \"disabled\",\"name\": {\"first\": \"Federico\",\"middle\": \"Kennedy\",\"last\": \"Armstrong\"},\"username\": \"Federico-Armstrong\",\"password\": \"aSzfK2V2XOgGO1H\",\"emails\": [\"Emelie98@example.com\",\"Jakayla1@gmail.com\"],\"location\": {\"street\": \"7213 Hagenes Cape\",\"city\": \"Kettering\",\"state\": \"North Carolina\",\"country\": \"Cambodia\",\"zip\": \"31336\",\"coordinates\": {\"latitude\": 68.4004,\"longitude\": -81.5312}},\"website\": \"https://severe-assist.info/\",\"domain\": \"quick-witted-default.org\",\"job\": {\"title\": \"Direct Paradigm Consultant\",\"descriptor\": \"Global\",\"area\": \"Infrastructure\",\"type\": \"Engineer\",\"company\": \"Hayes and Sons\"},\"creditCard\": {\"number\": \"3529-0083-5264-1525\",\"cvv\": 557,\"issuer\": \"mastercard\"},\"uuid\": \"5b8c7a77-d3e1-4547-ab65-6236710bac18\",\"objectId\": \"6857c01cf688af1fc5c99132\"}]}",
    "count" : 100000
}
```

#### Response

The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    "result": "Ok"
}
```
Failed result returns 400 status code (bad request)/500 status code (internal server error)
```
{
   "error": "error message"
}
```

-----------------------------------------------------------------------------------------

### /getWorkers

#### Definition
Receives a GET request for the pulling all the workers.

Receives no parameters and no body.

#### Example

URL 
```
http://localhost:8888/getWorkers
```
json body
```
no json body
```
#### Response

The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    "result": {
        "id": 1,
        "name": "Name"
    }
}
```
Failed result returns 500 status code (internal server error)
```
{
    "error": "error message"
}
```

-----------------------------------------------------------------------------------------

### /getRunningWorkers

#### Definition
Receives a GET request for the pulling all the running workers.

Receives no parameters and no body.

#### Example

URL 
```
http://localhost:8888/getRunningWorkers
```
json body
```
no json body
```
#### Response

The endpoint returns a status code and a json string.

Successed result returns 200 status code (OK) and 
```
{
    "result": {
        "id": 1,
        "name": "Name"
    }
}
```

-----------------------------------------------------------------------------------------

### /logsWs

#### Definition
Receives a GET request and establishes a realtime web sockets connection. 
After the connections established, the connection receives a message with a json body which contains the following fielsd:

|Name|data type|descriprion|
|:---|:---:|---|
| workerId | Integer | A running worker indentifire |

#### Example

URL 
```
ws://localhost:8888/logsWs
```
{
    "workerId": 1
}
#### Response

Web socket returns a message which contains an objects.

Message
```
{
    "Log": "Request 1 of TestWorker was success || Total amount of records have been sent 100 of 1000", 
    "RequestDurationTime": 0.01,
    "MedianRequestDurationTime": 0.01,
    "SuccessAttemptsCount": 10,
    "FailedAttemotsCount": 0
}
```

-----------------------------------------------------------------------------------------
