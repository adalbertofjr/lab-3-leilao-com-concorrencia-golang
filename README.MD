# Sistema de Leilão com Concorrência em Go

API REST para gerenciamento de leilões com processamento assíncrono de lances e fechamento automático de leilões.

## 🚀 Tecnologias

- Go 1.20
- MongoDB
- Gin Framework
- Docker & Docker Compose

## 📋 Pré-requisitos

- Docker e Docker Compose instalados
- Go 1.20+ (para desenvolvimento local)

## 🐳 Executando com Docker Compose (Recomendado)

### 1. Subir a aplicação

```bash
docker-compose up -d
```

Isso iniciará:
- **API** na porta `8080`
- **MongoDB** na porta `27017`
- **Mongo Express** (UI do MongoDB) na porta `8081`

### 2. Acessar o Mongo Express

Abra no navegador: http://localhost:8081

- **Usuário:** admin
- **Senha:** admin

### 3. Cadastrar usuários manualmente

⚠️ **IMPORTANTE:** Os usuários devem ser cadastrados manualmente no MongoDB antes de criar lances.

1. Acesse http://localhost:8081
2. Entre no database `auctions`
3. Acesse a collection `users`
4. Clique em **"New Document"**
5. Insira um documento no formato:

```json
{
  "_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva"
}
```

> **Dica:** Use um gerador de UUID online para criar IDs válidos.

### 4. Testar a API

Use os arquivos na pasta `api/` com a extensão REST Client do VS Code:

- [auction.http](api/auction.http) - Testes de leilões
- [bid.http](api/bid.http) - Testes de lances
- [user.http](api/user.http) - Testes de usuários

Ou use curl/Postman com os exemplos abaixo.

## 💻 Executando Localmente (Sem Docker)

### 1. Subir apenas o MongoDB

```bash
docker-compose up -d mongodb mongo-express
```

### 2. Configurar variáveis de ambiente

Edite o arquivo `cmd/auction/.env` se necessário.

### 3. Executar a aplicação

```bash
cd cmd/auction
go run main.go
```

A API estará disponível em: http://localhost:8080

## 🧪 Testes Automatizados

### Executar todos os testes

```bash
# Certifique-se que o MongoDB está rodando
docker-compose up -d mongodb

# Execute os testes
go test -v ./internal/infra/database/auction/
```

### Teste de fechamento automático de leilões

O teste `TestCreateAuction_AutoClose` valida que:
- Leilões são criados com status `Active`
- Após o intervalo configurado (`AUCTION_INTERVAL`), o status muda para `Completed`
- A goroutine de fechamento automático funciona corretamente

## 📡 Endpoints da API

### Leilões

```bash
# Criar leilão
POST http://localhost:8080/auction
Content-Type: application/json

{
  "product_name": "Notebook Dell",
  "category": "Eletrônicos",
  "description": "Notebook Dell Inspiron 15, i7, 16GB RAM",
  "condition": 1
}

# Listar leilões
GET http://localhost:8080/auction

# Buscar leilão por ID
GET http://localhost:8080/auction/{auctionId}

# Buscar lance vencedor
GET http://localhost:8080/auction/winner/{auctionId}
```

### Lances

```bash
# Criar lance
POST http://localhost:8080/bid
Content-Type: application/json

{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "auction_id": "08f770bf-929f-4f05-be85-83c692e065ce",
  "amount": 1500.00
}

# Listar lances de um leilão
GET http://localhost:8080/bid/{auctionId}
```

### Usuários

```bash
# Buscar usuário por ID
GET http://localhost:8080/user/{userId}
```

## ⚙️ Configurações

Edite `cmd/auction/.env`:

```env
# Intervalo para processar batch de lances
BATCH_INSERT_INTERVAL=20s

# Tamanho máximo do batch de lances
MAX_BATCH_SIZE=4

# Duração de um leilão antes de fechar automaticamente
AUCTION_INTERVAL=20s

# Conexão MongoDB
MONGODB_URL=mongodb://admin:admin@localhost:27017/auctions?authSource=admin
```

## 🎯 Condições de Produto

- `1` - Novo
- `2` - Usado
- `3` - Recondicionado

## 🛑 Parar a aplicação

```bash
docker-compose down
```

## 📝 Observações

- Lances são processados em **batch** para otimizar performance
- Leilões fecham **automaticamente** após o tempo configurado em `AUCTION_INTERVAL`
- Lances em leilões fechados são **rejeitados automaticamente**
- O sistema usa **goroutines** e **mutexes** para processamento concorrente seguro
