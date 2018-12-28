# UdaBayar Golang Backend

This is a Golang Backend for UdaBayar app. Here are the tech stack that this app used:

- Golang Language
- AWS SES as email service

## Installation

Make sure you already had golang installed on your PC, then run these steps:

- clone this repository

```
$ git clone https://github.com/Rekeningku/udabayar-go-api udabayar-go-api-di
```

- install all dependencies

```
$ go get
```

-  build this project.

```
$ make build
```

-  happy run with cli :)

```
$ ./udabayar
```

## Dev Database

- Database host etc available at environment

- All Migration file at directory model

- You can do migrate with command:

```
$ ./udabayar migrate database
```

- Seeding primary data and dummy data (if you needed):

```
$ ./udabayar seeder run
$ ./udabayar seeder dummy
```

## Development

- Run Development Mode with Udabayar CLI

```
$ ./udabayar serve development
```

## Deployment

- Make sure already done migrate database, seeding primary data and migrate email template

```
$ ./udabayar migrate database
$ ./udabayar seeder run
$ ./udabayar migrate email
```

- Deploy and run this app.

```
$ make build
```

## Roadmap

### Current Priorities

1) On Startup API (Generate categories, banks, user as super admin, and shipping_agents)
2) JWT user login and register
3) CRUDSS shop
4) CRUDSS storefront
5) CRUDSS product -> product_variant -> product_image (related)
6) CRUDSS user_bank
7) CRUDSS address
8) CRUDSS order and transaction (related) from customer
9) Event email
10) Third Party Bank API (from Defri)
11) Event Push Notification
12) Event SMS

Note: ✅ (Done), [Empty] (On Progress)

### ON STARTUP API

|       Startup API         | Status  |
|--------------------|:-----:|
| generateCategories   |      |
| generateBanks |      |
| generateUserSuperAdmin     |      |
| generateShippingAgents             |      |

### USER CREDENTIAL API

|       API         | Status  |
|--------------------|:-----:|
| login   |✅|
| register |✅|
| event-> sendEmailConfirmation |✅|
| validateToken |✅|
| forgotPassword |✅|
| changePassword |✅|

### SHOP API

|      API         | Status  |
|--------------------|:-----:|
| allShops   |✅|
| getShop |✅|
| createShop |✅|
| updateShop |✅|
| deleteShop |✅|

### STOREFRONT API

|       API         | Status  |
|--------------------|:-----:|
| allStorefronts   |✅|
| getStorefront |✅|
| createStorefront |✅|
| updateStorefront |✅|
| deleteStorefront |✅|

### PRODUCT API

|       API         | Status  |
|--------------------|:-----:|
| allProducts   |✅|
| getProduct |✅|
| uploadProductImage |✅|
| createProduct |✅|
| updateProduct |✅|
| deleteProduct |✅|

### USER_BANK API

|       API         | Status  |
|--------------------|:-----:|
| allUserBanks   |✅|
| getUserBank |✅|
| createUserBank |✅|
| updateUserBank |✅|
| deleteUserBank |✅|


### ADDRESS API

|       API         | Status  |
|--------------------|:-----:|
| allAddresses   |✅|
| getAddress |✅|
| createAddress |✅|
| updateAddress |✅|
| deleteAddress |✅|
| setMainAdress |      |
| getMainAdress |✅|

### ORDER & TRANSACTION API

|       API         | Status  |
|--------------------|:-----:|
| allTransactions   |✅|
| getTransaction |✅|
| createTransaction |✅|
| webhookTransaction |✅|
| event-> generateIdrPrefix |✅|
| trackTransaction |✅|

### THIRD PARTY BANK API

|       API         | Status  |
|--------------------|:-----:|
| generateIdrPrefix   |      |
| getIncomingTransfer |      |
