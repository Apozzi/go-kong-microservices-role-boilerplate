# Login Role API Backend
API feita em golang utilizando frameworks Gin e o a ORM do Gorm.


## Como rodar?
O backend roda utilizando Go 1.17.
Primeiramente é necessária instalação das dependencias na pasta do projeto.
```
go install
```
Após isso basta utilizar seguinte comando para rodar o backend:
```
go run main.go
```

### Deploy 
```
go deploy
```

## Swagger 

É possivel acessar o swagger do projeto em: http://localhost:8081/swagger/index.html#/

## Banco de dados

O projeto faz migração automatica.

Mas caso quiser adicionar as tabelas manualmente só acessar "database/addTables",
e caso quiser adicionar registros de exemplo "database/addData".

Banco de dados Postgres utilizado.
