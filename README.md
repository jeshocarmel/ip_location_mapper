
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



