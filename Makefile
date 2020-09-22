build:
	docker build -t jeshocarmel/ip_location_mapper:latest .

start:
ifdef apikey
	docker build -t jeshocarmel/ip_location_mapper:latest .
	docker run -it --env IPSTACK_API_KEY=$(apikey) -p 8080:8080  jeshocarmel/ip_location_mapper
else
	echo "Please mention apikey before you start"
endif

push:
	docker push jeshocarmel/ip_location_mapper	





