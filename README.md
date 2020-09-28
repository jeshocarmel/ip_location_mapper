
## IP Location Tracker

An app to trace the location of an IP address. The IP address is provided as input and the GPS coordinates for that IP is shown in a map.

**Flow:**
- API requests are made to [ipstack.com](https://ipstack.com/) to get the coordinates.
- The results are stored in a redis cache with TTL 24 hours.

**Architecture**

![Image of Yaktocat](https://raw.githubusercontent.com/jeshocarmel/ip_location_mapper/master/architecture.png)


### Prerequisites

- an api key from [ipstack](https://ipstack.com/) (free with 10,000 api requests per month )

### To start on a development machine

```bash
    docker-compose up --build
```

### To build the project and push to dockerhub

```bash
    make push
```

> Note: kindly go through the Makefile before performing a push. The dockerhub repository in Makefile points to **jeshocarmel/ip_location_mapper**


### To start the project on minikube

**1. start minikube**
```bash
minikube start
```

**2. create secret**
```
kubectl apply -f app-secret.yaml
kubectl describe secrets/app-secret
```

**3. add helm repo to install redis cluster**
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/redis --values minikube_files/values-minikube.yml
```

**4. create configmap**
```
kubectl apply -f app-configmap.yaml
kubectl describe configmap/app-configmap
```

**5. create deployment and service**
```
kubectl apply -f app-deployment.yaml
```

**6. view everything you have created**
```
kubectl get all
```

**7. expose service via minikube**
```
minikube service go-app-service
```

###### Additional Info

Checking logs of pods

- to view last 20 lines

```bash
kubectl logs --tail=20 deployment/go-app
```

- to stream logs

```bash
kubectl logs -f deployment/go-app
```

To delete everything (pods + deployment + service + statefulset) initiate the following commands
```
kubectl delete deployment/go-app
kubectl delete service/go-app-service
helm delete my-release
minikube delete
```