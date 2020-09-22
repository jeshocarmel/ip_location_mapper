### to start a docker container running redis on your local

docker build -t jeshocarmel/ip_location_mapper:latest .
docker run -it --env IPSTACK_API_KEY=<IPSTACK_API_KEY> -p 8080:8080 jeshocarmel/ip_location_mapper
docker push jeshocarmel/ip_location_mapper

docker run -p 6379:6379 --name some-redis -d redis --requirepass "YOUR_PASSWORD_HERE"
