# Notification Center
![CI status](https://img.shields.io/badge/build-passing-brightgreen.svg)

The notification center is a Golang microsystem to receive a large number of requests, send each of the payloads to Google Firebase and create persistence with MongoDB

## Run It
`$ go run *.go`

### Requirements
[mgo.V2](http://gopkg.in/mgo.v2)

```shell
$ go get gopkg.in/mgo.v2
```

[ab - Apache HTTP server benchmarking tool
](https://httpd.apache.org/docs/2.4/programs/ab.html) 

```shell
$ sudo apt-get install ab
```

## Setup

Following the good practices of the twelve-factor app, [chapter III - config](https://12factor.net/config)

```shell
$ export FIREBASE_KEY = {FirebaseKey}
```

## Tests
```shell
$ ab -p payload_test.data -T application/json -c 100 -n 20000 http://localhost:8080/push
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)