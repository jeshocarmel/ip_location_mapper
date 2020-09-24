
## IP Location Tracker

An app to trace the location of an IP address. The IP address is provided as input and the GPS coordinates for that IP is shown in a map.

**Flow:**
- API Requests are made to [ipstack.com](https://ipstack.com/) to get the coordinates.
- The results are stored in a redis cache with TTL 24 hours.

**Architecture**

![Image of Yaktocat](https://raw.githubusercontent.com/jeshocarmel/ip_location_mapper/master/architecture.png)


### Prerequisites

- an api key from [ipstack](https://ipstack.com/). free tier available with 10,000 api requests per month.

### To start the app on docker client

```bash
     make start apikey=<IPSTACK_API_KEY>
  ```

If you have a redis cluster/single node running anywhere, then use the below command.

```bash
    make start apikey=<IPSTACK_API_KEY> redis-url=<REDIS_HOST:PORT> redis-password=<REDIS_PASSWORD>
```

### To build the project and push to dockerhub

```bash
    make push
```

> Note: kindly view the Makefile before performing a push. The dockerhub repository in Makefile points to **jeshocarmel/ip_location_mapper**


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

