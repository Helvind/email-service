# email-service
A small service for sending e-mails via different providers

## Design
The service is designed with an API defined by the protocol buffer defined in `proto/email.proto`. Here we have the gRPC endpoint for sending e-mail, and message types for requests and responses. This gives is a programming language agnostict specification that we can easy use to generate code for new clients.

The structure of the service and the interface to the underlying e-mail providers is specified in `pkg/email/email.go`, while the implementations of the specific providers and the handling of sending is in `pkg/email/handler.go`. 

### Monitoring
The services is setup with a Prometheus scrap endpoint for metrics, to monitor traffic... TODO actual metrics

### Scaling
In `manifests/email-service.yaml` are definitions of Kubernetes resources, to deploy the service. This creates 2 replicas of the service that are loadbalanced through "TBD", and will automatically scale on CPU load up to 12 replicas through a Kubernetes HorizontalPodAutoscaler

## Running locally
The e-mail service can be started with:
```
go run cmd/service/service.go
```

When the service is running, it can be tested with the client:
```
go run cmd/client/client.go
```

### Testing
```
go test ./...
```
