
# NATS HTTP Adapter

### That is currently in Alpha and lack of a everything you need, containing only what i need in current state

This software is made to make proxy API to send NATS messages through HTTP for testing purposes.

The main idea is to proxy HTTP requests to NATS to existing API testing tools (i.e. Postman) might be used to test code which receives data through NATS


## Installation

First install command

```
go get github.com/paulvasilenko/nats-http-adapter
```

After create `config.yaml` file with following structure

```
NATS:
  Endpoint: gnats://endpoint:4222
  RequestTimeout: 500ms # Default is 1s
HTTP:
  Port: 85 # Default is 80
```

And run command

```
nats-http-adapter -c path/to/config.yaml
```


## Usage

Requests should be sent to endpoint `127.0.0.1:85/nats` if you run it locally. Change address from localhost to yours if necessary.

### Request

Example

```
{
    "subject":"ms_providers",
    "type":"req",
    "data": "any data"
}
```

List of params:
*  _subject_ - is a subject to which message is sent
* _type_ - is a type of nats publish. Available types:
  * _req_ - stands for *request*. Expects response, fails after timeout
  * _pub_ - stands for *publish*. Would just publish message and then return 204 No Content
* _data_ - raw text of data. Might be JSON or whatever you like.

### Response

#### OK

OK response is returned only for `req` message type. `data` contains and answer received through NATS.

If data is returned as JSON it would be encoded as JSON, otherwise as a raw string

```
{
    "data": "json object or string representation of data"
}
```

`pub` message type returns `204 No Content` always.

#### Error

Error response is returned only in case of adapter internal error, not receiver error

`400 Bad Request` returned if some params are missing

## TODO

* Only human-readable raw text is available now, support of protobuf should be implemented
* Implement other message types (RequestPublish, deliver etc)
* Wrap it into Docker container