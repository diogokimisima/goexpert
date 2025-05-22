# Weather System by CEP

Este projeto consiste em dois serviços que trabalham juntos para fornecer informações climáticas com base em um CEP brasileiro:

- **Serviço A**: API de entrada que valida o formato do CEP e encaminha a solicitação para o Serviço B
- **Serviço B**: Serviço de orquestração que busca os dados de localização e temperatura

O sistema utiliza OpenTelemetry para rastreamento distribuído e Zipkin para visualização de traces.

## Pré-requisitos

- Docker e Docker Compose
- Go 1.18 ou superior (para desenvolvimento)

## Como executar

1. Clone o repositório:

2. Inicie os serviços com Docker Compose:
    - docker-compose up -d

3. Verifique se todos os serviços estão em execução:
    - docker-compose ps

4. Teste o sistema enviando uma solicitação para o Serviço A:
    - curl -X POST http://localhost:8080/weather
    -H "Content-Type: application/json"
    -d '{"cep": "01001000"}'

## Endpoints

### Serviço A - Entrada de CEP (porta 8080)

- **POST /weather**
  - Corpo da requisição: `{"cep": "01001000"}`
  - Resposta de sucesso: Dados de temperatura e cidade do Serviço B
  - Resposta de erro: Código 422 para CEP inválido

### Serviço B - Processamento de dados climáticos (porta 8081)

- **POST /weather**
  - Corpo da requisição: `{"zipcode": "01001000"}`
  - Resposta de sucesso: `{"city": "São Paulo", "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65}`
  - Resposta de erro: Código 422 ou 404 para erros de CEP

## Visualização de traces

O Zipkin está disponível em: http://localhost:9411

## Serviços

- **service-a**: Serviço de entrada (porta 8080)
- **service-b**: Serviço de processamento climático (porta 8081)
- **otel-collector**: Coletor OpenTelemetry
- **zipkin**: Interface para visualização de traces

Passo 7: Configuração dos módulos Go
7.1 Para o Serviço A, você precisa configurar as dependências no arquivo go.mod:
    - cd service-a
    - go mod tidy

7.2 Para o Serviço B, atualize as dependências se necessário:
    -   cd service-b
        go get go.opentelemetry.io/otel
        go get go.opentelemetry.io/otel/sdk
        go get go.opentelemetry.io/otel/trace
        go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
        go mod tidy

Passo 8: Construção e teste dos serviços
Execute o Docker Compose para construir e executar todos os serviços:
    - docker-compose up --build

Após os serviços estarem em execução, teste o sistema enviando uma solicitação para o Serviço A:

    curl -X POST http://localhost:8080/weather \
    -H "Content-Type: application/json" \
    -d '{"cep": "01001000"}'