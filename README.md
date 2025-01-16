# Login Role API Backend
API feita em golang utilizando frameworks Gin e o a ORM do Gorm.


## Como rodar?
O backend roda utilizando Go 1.17.
Primeiramente é necessária instalação das dependencias na pasta do projeto.

vale lembrar que temos 3 projetos diferentes.

item-microservice
role-microservice
user-microservice

```
go install
```
Após isso basta utilizar seguinte comando para rodar o backend:
```
go run main.go
```

## Docker

Execute os seguintes comando na pasta de cada seguinte serviço

```
docker build -t item-microservice .
```

```
docker build -t role-microservice .
```

```
docker build -t user-microservice .
```
para gerar a imagens do docker.

E para rodar as imagens do docker respectivamente

```
docker run -p 8083:8083 item-microservice
```

```
docker run -p 8082:8082 role-microservice
```

```
docker run -p 8081:8081 user-microservice
```

ou utilizar ``` docker-compose up --build ``` para utilizar o docker-compose
