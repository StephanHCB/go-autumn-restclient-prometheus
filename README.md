# go-autumn-restclient-prometheus

Prometheus instrumentation functions for [go-autumn-restclient](https://github.com/StephanHCB/go-autumn-restclient).

## About go-autumn

A collection of libraries for [enterprise microservices](https://github.com/StephanHCB/go-mailer-service/blob/master/README.md) in golang that
- is heavily inspired by Spring Boot / Spring Cloud
- is very opinionated
- names modules by what they do
- unlike Spring Boot avoids certain types of auto-magical behaviour
- is not a library monolith, that is every part only depends on the api parts of the other components
  at most, and the api parts do not add any dependencies.  

Fall is my favourite season, so I'm calling it go-autumn.

## About go-autumn-restclient

It's a rest client that also supports x-www-form-urlencoded.

## About go-autumn-restclient-prometheus

Implements instrumentation callbacks that use [prometheus/client_golang](https://github.com/prometheus/client_golang).

## Usage

Use the provided callbacks while constructing your rest client stack.
