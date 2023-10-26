
# MNC Test


#### Example Usage

```go
	go run cmd/main.go
```

#### Copy whole .env.example and create a new .env

#### Before running this app, make sure you have already installed Redis and connected to it first

#### Paseto expire token (Hour Format) ex :
```go
jsonToken := paseto.JSONToken{
		Issuer:     "Sora Project",
		Subject:    "Abrakadabra",
		Expiration: expire,
		IssuedAt:   now,
	}
```
#### Endpoints Available

```http
POST   /auth/register  //register new customer         
POST   /auth/login     //customer log in          
POST   /auth/logout    //coustomer log out         
```

#### Endpoints (Middleware Protected)
```http
GET    /app/merchants/list   //catch available merchants on database (middleware area)
POST   /app/merchants/create   //create new mercahnts 
POST   /app/transaction/create //make a new payment to merchant
POST   /app/transaction/list    //Showing transaction list (that has been made)
POST   /app/transfer/create/account //make a new trransfer 
POST   /app/transfer/list/income //Showing income history
POST   /app/transfer/list/outcome //Showing outcome history
```

| Description                |
| :------------------------- |
| **Required**. Paseto / JWT Token |

#### Activity Logs & HTTP Request Incoming
Recorded and saved to activity.log & request.log file
#### example :
Customer has been logged in :
```
{"level":"info","ts":1698301464.8801463,"caller":"usecase/user-credentials-usecase.go:113","msg":"Customer has been logged in","Customer Email":"awd@gmail.com"}
```
Customer has made a new payment :
```
{"level":"info","ts":1698258809.2386837,"caller":"usecase/transaction_usecase.go:35","msg":"A payment has been made","customerID":3,"merchantID":1,"amount":5000}
```
Customer has successfully transferred money: :
```
{"level":"info","ts":1698291841.860478,"caller":"usecase/transfer_usecase.go:43","msg":"Request transfer money has been initiated","senderAcountNumber":"12481257","receiverAccountNumber":"12371246","amount":10000}
```
Requested to create a new merchant :
```
{"level":"info","ts":1698319011.2144413,"caller":"controller/merchant-controller.go:35","msg":"New merchant has been created","Merchant Name":"haji barokah"}
```

#### Logout Method
```http
POST   /auth/logout    //coustomer log out  
```

1.User logs in: When a user logs in, they provide their credentials (typically an email and password). The server verifies these credentials, and if they are correct, it generates a unique token for the user. This token is sent back to the user and also stored in Redis with an expiration time of 24 hours.

2.User logs out: When a user wants to log out, they send a request to the logout endpoint with their email in the JSON body. The server then deletes the corresponding entry in Redis.

3.Post-logout access attempts: If a user tries to access a protected endpoint after logging out, their token will no longer be found in Redis when the middleware checks it. Therefore, the access attempt will fail, effectively logging out the user.

#### Postman Documentation
https://documenter.getpostman.com/view/29723627/2s9YRFTov9

#### Postman Collection
https://www.postman.com/crimson-crater-616314/workspace/mnc-test/collection/29723627-aed2b8b7-fb73-42b9-be9e-9457d38871ba?action=share&creator=29723627

#### Unit test result
![alt text](https://i.ibb.co/cFxrvLL/unit-test-res.jpg)

#### Stack
- Gin (Web Framework)
- Paseto (Security concern)
- Zap (Logging)
- Redis (Caching)
