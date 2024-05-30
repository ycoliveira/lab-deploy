# goexpert-lab-cloud-run

## Commands to start the application locally

1. Just run the command below to start 

```sh
docker compose up -d
```

2. The server will be available at the url below:

```sh
http://localhost:8080
```

3. The main route is accessed by informing a valid CEP as a path parameter. 
http://localhost:8080/weather/{cep}

  For example:
```sh
http://localhost:8080/weather/01153000
```

## Informations regarding the Google Cloud run deploy

The web application deployed using Google Cloud Run can be accessed using the address https://lab-cloud-run-7s5babtwoa-uc.a.run.app/{cep}

For example:
```sh
https://lab-cloud-run-7s5babtwoa-uc.a.run.app/01153000
```
