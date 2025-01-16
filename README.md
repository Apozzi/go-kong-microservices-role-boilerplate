# Login Role API Backend
Microserviços feitos em golang utilizando frameworks Gin e o a ORM do Gorm e Kong como Api Gateway.


## Como rodar?
O backend roda utilizando Go 1.17, é importante que tenha Postgres instalado e rodando na maquina.

Primeiramente é necessária instalação das dependencias na pasta de cada projeto, vale lembrar que temos 3 pastas com microserviços diferentes.

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

Ao rodar cada microserviço ele fará migrate na base criando todas as tabelas no Postgres, também é possivel fazer criação de tabela e registros através da pasta `database`.

## Usando Kong e Konga, como Api Gateway.

Dentro das pasta `docker-kong` criar imagem docker com commando ```docker network create kong-net```, execute os seguintes comandos respectivamente:
- ```docker-compose up -d db``` esperar iniciar o banco.
- ```docker-compose ps```
- ```docker-compose up -d ```

Após isso acessar `localhost:1337` e criar usuario e senha (a que eu usei foram username:admin, password:ADMINADMIN).

