
## IP Location Tracker

### To start the app on development server

```bash
     make start apikey=<IPSTACK_API_KEY>
  ```

    If you have a redis cluster running anywhere, do mention it as well here

```bash
    make start apikey=<IPSTACK_API_KEY> redis-url=<REDIS_HOST:PORT> redis-password=<REDIS_PASSWORD>
```

### To build the project and push to dockerhub

```bash
    make push
```


