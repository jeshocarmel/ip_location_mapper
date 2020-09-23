
## IP Location Tracker

### Prerequisites

- an api key from [ipstack](https://ipstack.com/). ipstack provides 10,000 api requests free per month.

### To start the app on docker client

```bash
     make start apikey=<IPSTACK_API_KEY>
  ```

If you have a redis cluster running anywhere, then use the below command.

```bash
    make start apikey=<IPSTACK_API_KEY> redis-url=<REDIS_HOST:PORT> redis-password=<REDIS_PASSWORD>
```

### To build the project and push to dockerhub

```bash
    make push
```

> Note: kindly check the Makefile before performing a push. The dockerhub repository in Makefile points to **jeshocarmel/ip_location_mapper**


### To start the project on minikube

```bash
minikube start
kubectl apply -f app-secret.yaml
kubectl apply -f app-deployment.yaml
kubectl apply -f app-service.yaml 
minikube service go-app-service
```

To check the logs of the pods 

- to view last 20 lines

```bash
kubectl logs --tail=20 deployment/go-app
```

- to stream logs

```bash
kubectl logs -f deployment/go-app
```

