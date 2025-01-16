# Login Role Microservices Backend
Microserviços feitos em golang utilizando frameworks Gin e o a ORM do Gorm e Kong como Api Gateway.

Os serviços tem modelos de usuário e permissões por papeis.


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

Após isso acessar `localhost:1337` aonde fica o Konga e criar usuario e senha (a que eu usei foram username:admin, password:ADMINADMIN).

Depois disso é só criar seguinte conexão com o Kong.
![Conexao](https://github.com/user-attachments/assets/a896a1f8-b5e2-4d2a-b79c-0d10f60e3d92)

Na aba serviços criar cada respectivo serviço `rolemanager`, `itemmanager`, `usermanager`, da seguinte forma (vale lembrar que no host utilizaremos ip da nossa maquina local).

![Capturar](https://github.com/user-attachments/assets/de89f2af-5118-46d0-83de-fabb29c145a5)

E após isso vamos e `routes` e criamos uma rota também para cada serviço.

![Capturar2](https://github.com/user-attachments/assets/2da151c1-b0b1-4e45-bc50-c7cc82759d50)

Com isso nossos serviços estarão todos rodando através da mesma porta no `localhost:8000` nas rotas `http://localhost:8000/usermanager`, `http://localhost:8000/itemmanager`, `http://localhost:8000/rolemanager`, e 
podemos acessar swagger de cada um respectivamente: http://localhost:8000/usermanager/swagger/index.html#/ (é possivel editar as rotas no swagger para aparecer da forma correta após config do Kong)




