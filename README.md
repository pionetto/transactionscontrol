<h1 align="center">:file_cabinet: Transactions Control - RESERVA</h1>

## :memo: Descrição
Este projeto é um software para controle de transações de cartão de alimentação, refeição e dinheiro.
Ele .

## :books: Funcionalidades
* <b>Funcionalidades</b>
</br>
Esta é uma API desenvolvida na linguagem [Golang](https://go.dev/).

Ela realiza:

* Cadastra uma conta com 3 tipos de saldo;
* Registra transações entre contas (autorizadas ou negadas);
* Atualiza saldos de contas;
* Visualiza listagem das contas;
* Visualiza listagem das transações;
* Busca um estabelecimento pelo nome;
* Busca todos os resultados por conta;


O usuário realiza o login via middleware, através de token de autentitação 
via JWT.

Abaixo estão as rotas criadas:

```go

func GetRouter(app *app.App) *mux.Router {

	// middleware compartilhado em todas as rotas da API
	common := negroni.New(
		negroni.NewLogger(),
	)

	// criando roteador base
	router := mux.NewRouter()

	// rota de login
	loginRoutes := mux.NewRouter()
	router.Path("/login").Handler(common.With(
		negroni.Wrap(loginRoutes),
	))
	logins := loginRoutes.Path("/login").Subrouter()
	logins.Methods("POST").HandlerFunc(login.HandlerLogin(app))

	// rota de accounts
	accountsRoutes := mux.NewRouter()
	router.Path("/accounts").Handler(common.With(
		negroni.Wrap(accountsRoutes),
	))
	accounts := accountsRoutes.Path("/accounts").Subrouter()
	accounts.Methods("GET").HandlerFunc(account.ListAccounts(app))
	accounts.Methods("POST").HandlerFunc(account.PostAccount(app))

	// rota de balance
	balanceRoutes := mux.NewRouter()
	router.Path("/accounts/{id}/balance").Handler(common.With(
		negroni.Wrap(balanceRoutes),
	))
	balances := balanceRoutes.Path("/accounts/{id}/balance").Subrouter()
	balances.Methods("GET").HandlerFunc(account.BalanceAccount(app))

	// rota de transações (transactions)
	transactionsRoutes := mux.NewRouter()
	router.Path("/transactions").Handler(common.With(
		negroni.Wrap(transactionsRoutes),
	))
	transactions := transactionsRoutes.Path("/transactions").Subrouter()
	transactions.Methods("GET").HandlerFunc(transaction.ListTransactions(app))
	transactions.Methods("POST").HandlerFunc(transaction.PostTransactions(app))

	// rota de estabelecimentos
	merchantsRoutes := mux.NewRouter()
	router.Path("/merchants/{merchant}").Handler(common.With(
		negroni.Wrap(merchantsRoutes),
	))
	merchants := merchantsRoutes.Path("/merchants/{merchant}").Subrouter()
	merchants.Methods("GET").HandlerFunc(merchant.ListMerchants(app))

	return router
}

```
E para acessa-las, podemos utilizar uma ferramenta chamada Insomnia(clique [aqui](https://insomnia.rest/download) para baixar),
para testar se os endpoints estão funcionando corretamente.
Ela serve para que possamos fazer requisições (GET, POST, PUT, DELETE..).

## Sobre a API
* Todos os caminhos da API poderão ser acessados a partir do link http://localhost:8080;
* As respostas das requisições feitas a API são em formato JSON;

## Accounts (Contas)
</br>

**Criar uma conta**
</br>

**Método:** POST
</br>

**Endpoint:** http://localhost:8080/accounts
</br>

**Objeto JSON a ser enviado:**

```JSON
{
	"CPF": "11111111111",
	"secret": "123456",
	"amount_food": 1000.00,
	"amount_meal": 1000.00,
	"amount_cash": 1000.00
}
```

**Listar contas**
</br>

**Método:** GET
</br>

**Endpoint:** http://localhost:8080/accounts
</br>
</br>


## Transactions (Transações)
</br>

**Criar transações**
</br>

**Método:** POST
</br>

**Endpoint:** http://localhost:8080/transactions
</br>

**Objeto JSON a ser enviado:**

```JSON
{
	"accounttocredit_id": 1,
	"amount": 1000.00,
	"merchant": "Mercantil Dourado",
	"mcc": ""
}
```
</br>
Neste caso o campo "mcc" do JSON deve receber como parâmetro o código correspondente ao da transação,
sendo eles:
</br>

**FOOD** - 5411 ou 5412
**MEAL** - 5811 ou 5812
**CASH** - Campo vazio ou qualquer número diferente dos anteriores.

</br>

**Listar transações** </br>
</br>

**Método:** GET
</br>

**Endpoint:** http://localhost:8080/transactions
</br>


# Merchants (Estabelecimentos)

**Listar estabelecimentos** </br>
</br>

**Método:** GET 
</br>

**Endpoint:** http://localhost:8080/merchants/
</br>

**Listar um estabelecimento pelo nome**
</br>

**Método:** GET
</br>

**Endpoint:** http://localhost:8080/merchants/{merchant}

# Login

**Criar token de autenticação**
</br>

**Método:** POST
</br>

**Endpoint:** http://localhost:8080/login
</br>

**Objeto JSON a ser enviado:**

```JSON
{
	"cpf": "11111111111",
	"secret": "123456"
}
```

Após isso, retornará um JSON nesse formato:

```JSON
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcGYiOiI0NDQ0NDQ0NDQ0NCIsImV4cCI6MTY5NDM2MzMwNn0.ZfBdWafcs73NdM4rs4_duvggE0BZ4MB9UM9vxj3_yv4"
}
```
</br>
Este token deve ser utilizado para Listar Transações e Listar Estabelecimentos.


## :wrench: Tecnologias utilizadas
* [Golang](https://go.dev/);
* [Gorm](gorm.io/gorm);
* [Golang-JWT](https://github.com/golang-jwt/jwt);
* [Gorilla Mux](github.com/gorilla/mux);
* [PostgreSQL](https://go.dev/);
* [Github](https://go.dev/);
* [Visual Studio Code](https://go.dev/);
* [S.O. Linux Mint](https://go.dev/);
* [Insomnia](https://go.dev/) - Para alguns testes da API;

## :rocket: Rodando o projeto
São necessários alguns requisitos para rodar a aplicação.

* Ter o PostgreSQL instalado e configurado, tem um guia clicando [aqui](https://www.edivaldobrito.com.br/como-instalar-o-postgresql-no-ubuntu-20-04-lts-e-derivados/).
* Ter instalado o Golang.

Para instalar o golang, faça:

```
sudo apt install -y golang
```

Crie um diretório (mkdir) em seu lugar de preferência.

```
sudo apt install -y golang
```
Para rodar o repositório é necessário clonar o mesmo, dar o seguinte comando para iniciar o projeto:

```
git clone https://github.com/pionetto/transactionscontrol.git
```

É necessário configrar as variáveis de ambiente criando um arquivo `.env`:

```
BUILD_TARGET="development"
DEBUG_MODE="false"
TOKEN_KEY="gophers"
SERVER_ADDRESS="8080"
POSTGRES_PASSWORD="postgres"
POSTGRES_USER="postgres"
POSTGRES_PORT="5432"
POSTGRES_HOST="localhost"
POSTGRES_DB="caju"
```

Em seguida, para instalar as dependências do projeto e executar a aplicação,
acesse a raiz da pasta do projeto e digite:

```
go run main.go
```

O que ainda está sendo implementado?

* Banco de dados hospedado na AWS-RDS
* API hospedada na AWS-EC2
* Frontend em ReactJS. https://hygya-interface.vercel.app/login
* Repositório do Frontend: https://github.com/pionetto/hygya-interface

## :soon: Implementação futura
* O que será implementado na próxima sprint?

## :handshake: Colaboradores
<table>
  <tr>
    <td align="center">
      <a href="http://github.com/pionetto">
        <img src="https://avatars.githubusercontent.com/u/5672555?v=4" width="100px;" alt="Foto de Pio Netto no GitHub"/><br>
        <sub>
          <b>@pionetto</b>
        </sub>
      </a>
    </td>
  </tr>
</table>

## :dart: Status do projeto