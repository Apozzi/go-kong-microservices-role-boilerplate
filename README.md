# Login Role Microservices Backend
Microserviços feitos em golang utilizando frameworks Gin e o a ORM do Gorm e Kong como Api Gateway, e com RabbitMq.

Os serviços tem modelos de usuário e permissões por papeis.

## Arquitetura

![Capturar](https://github.com/user-attachments/assets/0cf6b2f4-b8ac-4fe6-80c5-2d55ca12a0d9)

## Como rodar? (Windows)
O backend roda utilizando Go 1.17, é importante que tenha Postgres instalado e rodando na maquina.

Primeiramente devemos ter broker executando na nossa maquina, com docker instalado, podemos rodar através do comando docker.:

`docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4.0-management` 

Com isso verificamos se rabbitMq está rodando e se quisermos podemos acessar `http://localhost:15672/#/` (User:guest, senha:guest)

Após isso necessária instalação das dependencias na pasta de cada projeto, vale lembrar que temos 4 pastas com microserviços diferentes.

`item-microservice`

`role-microservice`

`user-microservice`

`email-microservice`

Rodar o seguinte comandos em todos os projetos.

```
go install
```
Após isso basta utilizar seguinte comando para rodar o backend em cada pasta de cada microserviço:
```
go run cmd/main.go
```

Com isso teremos os seguintes microserviços rodando nas seguintes portas.

`item-microservice` = `localhost:8083`

`role-microservice` = `localhost:8082`

`user-microservice` = `localhost:8081`

O `email-microservice` não roda nenhuma porta e só serve para se comunicar com outros serviços através do rabbitMq.

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

<!---
 restruturar https://devopsian.net/p/how-to-structure-a-go-project-start-simple-refactor-later/
--!>
