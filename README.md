# loadbalancer-controller
Loadbalancer operator working on Kubernetes

# Resources

## Loadbalancer

Loadbalancer is a resource that define abstracted loadbalancer and healthcheck. 
Real loadbalancers are defined by `backend` resource.

There are 2 backend resources, `AWSBackend` and `ExternalBackend`

### AWSBackend

AWSBackend is a resource that is made by elastic loadbalancing with AWS.

See this example.

### ExternalBackend

ExternalBackend is a resource that can define any impremented loadbalancer.

See this example.

## LBSets

LBSet is a resource that join multiple `Loadbalancer` as a GSLB (Global Server Load Balancing)

There is a `Route53` Resource as a backend implementation.
