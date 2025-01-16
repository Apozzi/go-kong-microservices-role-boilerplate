# Login Role API Backend
API feita em golang utilizando frameworks Gin e o a ORM do Gorm.


## Como rodar?
O backend roda utilizando Go 1.17.
Primeiramente é necessária instalação das dependencias na pasta do projeto.

vale lembrar que temos 3 projetos diferentes.

`item-microservice`

`role-microservice`

`user-microservice`

Rodar o seguinte comandos em todos os projetos.

```
go install
```
Após isso basta utilizar seguinte comando para rodar o backend em cada pasta de cada microserviço:
```
go run main.go
```

Com isso teremos os seguintes microserviços rodando nas seguintes portas.

`item-microservice` = `localhost:8083`

`role-microservice` = `localhost:8082`

`user-microservice` = `localhost:8081`

## Usando Kong e Konga, como Api Gateway.
