# Client-Server API - Cotação de Dólar

Este repositório contém a solução para o desafio de implementar dois sistemas em Go: um servidor (`server.go`) e um cliente (`client.go`). O servidor consome uma API de câmbio, registra as cotações em um banco de dados e responde a requisições HTTP com o valor atual do câmbio. O cliente consulta o servidor, recebe a cotação e salva o valor em um arquivo de texto.

## Requisitos

- **server.go**: Deve realizar uma requisição HTTP para obter a cotação do dólar e registrar os dados no banco de dados SQLite.
- **client.go**: Deve solicitar a cotação ao servidor e salvar o valor do câmbio em um arquivo `cotacao.txt`.

## Funcionalidades

### server.go

- O servidor consome a API externa [AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL) para obter a cotação do dólar (USD-BRL).
- O servidor expõe o endpoint `/cotacao`, que retorna o valor do câmbio em formato JSON.
- O servidor usa **contexto** para limitar os tempos de execução:
  - Timeout máximo de 200ms para obter a cotação da API.
  - Timeout máximo de 10ms para salvar os dados no banco de dados.
- O servidor salva os dados da cotação no banco de dados SQLite.

### client.go

- O cliente envia uma requisição HTTP para o servidor e recebe a cotação atual do dólar.
- O cliente usa **contexto** para limitar o tempo de espera para 300ms.
- Após receber a cotação, o cliente salva o valor no arquivo `cotacao.txt` no formato: `Dólar: {valor}`.

## Estrutura do Projeto

```
├── client/
│   └── client.go           # Código do cliente para consultar o servidor e salvar a cotação
├── server/
│   └── server.go           # Código do servidor que fornece a cotação e a salva no banco de dados
├── main.go                 # Função principal para iniciar o servidor
├── test.http               # Teste de solicitação HTTP para testar a API
└── init.sql                # Script para inicializar o banco de dados SQLite
```

## Pré-requisitos

- Go (versão 1.18 ou superior)
- SQLite

## Como Usar

### 1. Inicializar o Banco de Dados SQLite

O arquivo `init.sql` contém o script para criar a tabela `exchange_rate` no banco de dados SQLite. O banco de dados será utilizado para armazenar as cotações recebidas.

### 2. Subir a Aplicação com Docker

Se você preferir usar Docker, é possível subir a aplicação utilizando Docker Compose. Para isso, siga os passos abaixo:

1. Certifique-se de ter o Docker e Docker Compose instalados em sua máquina.
2. No diretório raiz do projeto, execute o comando para construir as imagens e subir os containers:

```bash
docker-compose up --build
```

Este comando irá:

- Construir a imagem do Docker para o servidor e o cliente.
- Subir os containers do servidor e do banco de dados (caso o banco seja configurado para rodar dentro de um container).
- O servidor estará disponível na porta `8080` dentro do container.

### 3. Iniciar o Servidor

Caso não utilize Docker, você pode rodar o servidor diretamente:

1. Navegue até o diretório `server/`.
2. Execute o servidor com o comando:

```bash
go run server.go
```

### 4. Consultar a Cotação com o Cliente

O cliente fará uma requisição ao servidor para obter a cotação do dólar e salvá-la em um arquivo `cotacao.txt`.

1. Navegue até o diretório `client/`.
2. Execute o cliente com o comando:

```bash
go run client.go
```

O cliente salvará o valor do câmbio no arquivo `cotacao.txt` com o formato: `Dólar: {valor}`.

### 5. Testar a API Manualmente

Você pode testar o servidor manualmente utilizando o arquivo `test.http`. Esse arquivo contém uma requisição HTTP para testar o endpoint `/cotacao`.

A requisição para obter a cotação é:

```http
GET http://localhost:8080/cotacao
```

### 6. Banco de Dados SQLite

O banco de dados SQLite será utilizado para armazenar as cotações recebidas da API externa. A tabela `exchange_rate` conterá as seguintes colunas:

- `id`: UUID gerado automaticamente para cada registro.
- `code`: Código da moeda (ex: USD).
- `codein`: Código da moeda de destino (ex: BRL).
- `name`: Nome da moeda.
- `high`, `low`: Valores máximos e mínimos da cotação.
- `varBid`, `pctChange`: Variação do valor.
- `bid`, `ask`: Valores de compra e venda.
- `timestamp`: Timestamp da cotação.
- `create_date`: Data de criação do registro.

## Código e Fluxo de Execução

### `server.go`

1. O servidor faz uma solicitação para a API externa com um **contexto** que define um tempo máximo de 200ms para obter a cotação.
2. O servidor, em seguida, armazena as cotações no banco de dados SQLite com um **contexto** que define um tempo máximo de 10ms para persistir os dados.
3. O servidor responde à requisição com a cotação atual.

### `client.go`

1. O cliente envia uma requisição ao servidor para obter a cotação.
2. O cliente usa um **contexto** para garantir que o tempo de espera para a resposta do servidor não ultrapasse 300ms.
3. O cliente salva a cotação no arquivo `cotacao.txt`.

## Observações
Diferente do desafio original, o banco de dados utilizado foi o **MySQL**, devido a algumas limitações da máquina que impediram o uso do SQLite.

Além disso, os tempos de timeout foram ajustados, pois os valores propostos inicialmente para o desafio eram excessivamente baixos, o que poderia comprometer a performance e a estabilidade do sistema.