# Build stage
FROM golang:1.22-alpine AS builder

# Definir diretório de trabalho
WORKDIR /app

# Copiar os arquivos go.mod e go.sum (se existirem)
COPY go.mod .
COPY go.sum* .

# Download das dependências
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado do stage anterior
COPY --from=builder /app/main .

# Expor a porta que sua aplicação usa (ajuste conforme necessário)
EXPOSE 8082

# Executar a aplicação
CMD ["./main"]