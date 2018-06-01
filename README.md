# Fake MMS

Fake MMS is an ultra lightweight mocking server for the AWS [Marketplace Metering
Service](https://docs.aws.amazon.com/marketplacemetering/latest/APIReference/Welcome.html).


## Installation

```
go install github.com/kellyp/fake_mms
```

## Running

```
fake_sqs
```

## Development

```
go get github.com/kellyp/fake_mms
cd $GOPATH/src/github.com/kellyp/fake_mms
go build
```

## Docker

### Build

```
cd $GOPATH/src/github.com/kellyp/fake_mms
docker build -t your/tag_name .
```

Build the fake mms docker image and tag it appropriately

### Running

```
docker run -it -p 8080:8080 kellyp/fake_mms
```

bind to the host port 8080 and run a docker container
with fake mms
