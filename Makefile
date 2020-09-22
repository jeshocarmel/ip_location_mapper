build:
	docker build -t jeshocarmel/ip_location_mapper:latest .

start:
ifdef apikey
	make build
	docker run -it --env IPSTACK_API_KEY=$(apikey) --env REDIS_URL=$(redis-url) --env REDIS_PASSWORD=$(redis-password) -p 8080:8080  jeshocarmel/ip_location_mapper
else
	echo "Please mention apikey before you start"
endif

push:
	make build
	docker push jeshocarmel/ip_location_mapper	





