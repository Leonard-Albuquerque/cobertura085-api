# Cobertura085 API 🚚📍

API RESTful backend desenvolvida em **Go (Golang)** projetada para centralizar a camada de dados, regras de negócio e telemetria da plataforma de frete e cobertura por bairros em Fortaleza, desencaixando o banco de dados das Server Actions do frontend Next.js.

---

## 🎯 Proposta do Projeto

Esta API fornece endpoints de alta performance (<20ms) para:
- Gerenciamento de **Lojas / Estabelecimentos** e seus **Pontos de Retirada** (`PickupPoint`).
- Matriz de **Taxas de Entrega e Bairros** atendidos por loja em Fortaleza.
- Registro e agregação de **Telemetria de Buscas** em tempo real.
- Autenticação e gestão de usuários via **JWT**.

---

## 📐 Arquitetura do Sistema

A API segue os princípios de **Clean Architecture / Layered Design**:

```text
cobertura085-api/
├── cmd/
│   └── api/             # Ponto de entrada (main.go e registro de rotas Gin)
├── internal/
│   ├── config/          # Carregamento de variáveis de ambiente
│   ├── database/        # Conexão GORM + PostgreSQL
│   ├── dto/             # Data Transfer Objects (Request/Response contracts)
│   ├── handler/         # Controllers HTTP (Gin Handlers)
│   ├── middleware/      # Middlewares (Autenticação JWT, CORS, Logger)
│   ├── model/           # Entidades GORM (Mapeamento do PostgreSQL)
│   ├── repository/      # Camada de Persistência (Queries GORM)
│   └── service/         # Camada de Regras de Negócio e Sanitização
└── migrations/          # Migrações SQL do banco de dados
```

---

## 🚀 Implementações Realizadas

### 🏢 1. Módulo de Lojas (`/api/v1/stores`)
- `GET /api/v1/stores` — Lista todas as lojas cadastradas para a busca pública inicial.
- `GET /api/v1/stores/:id` — Retorna dados completos da loja (por slug ou ID) incluindo `pickupPoints`.
- `GET /api/v1/stores/qr/:token` — Resolve loja associada a um QR Token para redirecionamento.
- `PUT /api/v1/stores/:id/settings` — Atualização atômica de configurações da loja e sincronização da lista de pontos de retirada.
- `GET /api/v1/stores/:id/dashboard-stats` — Retorna estatísticas de bairros ativos/inativos e frete médio para o painel lojista.

### 📍 2. Módulo de Frete & Consultas Externas (`/api/v1/shipping`) [Fase 2]
- `POST /api/v1/shipping/lookup-cep` — Consulta de CEP via **ViaCEP**, verificação de frete na loja e gravação automática de telemetria.
- `POST /api/v1/shipping/lookup-address` — Geocodificação de endereço via **OpenStreetMap Nominatim**, extração de bairro e verificação de taxa.
- `POST /api/v1/shipping/lookup-coords` — Geocodificação reversa por latitude/longitude via **Nominatim**, identificação de bairro e taxa.
- `POST /api/v1/shipping/lookup-selected-address` — Processamento de endereço pré-selecionado no autocompletar.
- `GET /api/v1/shipping/address-suggestions` — Retorna até 5 sugestões de endereço em Fortaleza para o autocompletar em tempo real.

### 🏘️ 3. Módulo de Bairros e Taxas (`/api/v1/neighborhoods`)
- `GET /api/v1/stores/:id/neighborhoods` — Retorna a matriz de bairros e taxas da loja ordenada pelo nome oficial.
- `GET /api/v1/stores/:id/neighborhoods/check` — Consulta de frete por `baseNeighborhoodId`.
- `PATCH /api/v1/neighborhoods/:id` — Atualiza taxa, prazos, pedido mínimo e notas de um único bairro.
- `PATCH /api/v1/neighborhoods/bulk` — Atualização em massa de regras de entrega para múltiplos bairros.

### 📋 4. Módulo de Domínio & Utilitários (`/api/v1`)
- `GET /api/v1/lines-of-business` — Lista ramos de atuação ativos no sistema.
- `GET /api/v1/base-neighborhoods/by-name/:name` — Consulta normalizada de bairros de Fortaleza.
- `GET /api/v1/geojson/bairros-fortaleza` — Serve o GeoJSON com os limites territoriais dos bairros de Fortaleza para os mapas interativos Leaflet. [Fase 2]

### 📊 5. Módulo de Telemetria & Analytics (`/api/v1/telemetry`)
- `POST /api/v1/telemetry/search-events` — Grava log de consulta com IP anonimizado (SHA-256).
- `GET /api/v1/telemetry/summary` — Métricas globais de buscas, taxa de disponibilidade, latência média e top 10 bairros.
- `GET /api/v1/telemetry/logs` — Histórico recente de buscas para auditoria.
- `GET /api/v1/telemetry/stores/:storeId` — Métricas segregadas por empresa.

### 🔑 6. Módulo de Autenticação (`/api/v1/auth`)
- `POST /api/v1/auth/register` — Registro de novo usuário lojista.
- `POST /api/v1/auth/login` — Autenticação e emissão de JWT Access + Refresh Tokens.
- `POST /api/v1/auth/refresh` — Renovação do Access Token.
- `GET /api/v1/auth/me` — Dados do perfil do usuário autenticado.

---

## 🛠️ Como Executar

### Pré-requisitos
- Go 1.25+
- PostgreSQL (ou Neon DB)

### Passos
1. **Configurar Variáveis de Ambiente:**
   Crie o arquivo `.env` baseado no `.env.example`:
   ```env
   PORT=8080
   GIN_MODE=debug
   DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require
   JWT_SECRET=seu_secret_jwt_aqui
   ```

2. **Compilar / Executar a API:**
   ```bash
   go run ./cmd/api
   ```

3. **Executar Testes Unitários:**
   ```bash
   go test ./...
   ```
