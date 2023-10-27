
# MNC Test


#### Example Usage

```go
	go run main.go
```

#### Copy whole .env.example and create a new .env

#### Before running this app, make sure you have already installed Redis and connected to it first

#### Paseto expire token (Hour Format) on .env :
```
PASETO_EXP=5
```
#### Endpoints Available

```http
POST   /auth/register  //register new customer         
POST   /auth/login     //customer log in          
POST   /auth/logout    //coustomer log out         
```

#### Endpoints (Middleware Protected)
```http
GET    /app/merchants/list   //catch available merchants on database
POST   /app/merchants/create   //create new merchant 
POST   /app/transaction/create //make a new payment to merchant
POST   /app/transaction/list    //Showing transaction list (that has been made)
POST   /app/transfer/create/account //make a new transfer 
POST   /app/transfer/list/income //Showing income history
POST   /app/transfer/list/outcome //Showing outcome history
```

| Description                |
| :------------------------- |
| **Required**. Paseto Token |

#### Activity Logs & HTTP Request Incoming
Recorded and saved to activity.json & request.json file
#### example :
Customer has been logged in :
```
{"level":"info","ts":1698318833.313814,"caller":"controller/user-credential-controller.go:77","msg":"Customer has been logged in","Customer Email":"pall12@gmail.com"}
```
Customer has been logged out :
```
{"level":"info","ts":1698318779.6292548,"caller":"controller/user-credential-controller.go:107","msg":"Customer has been logged out","Customer Email":"pall12@gmail.com"}
```
Customer has made a new payment :
```
{"level":"info","ts":1698333614.1180406,"caller":"controller/transaction_controller.go:44","msg":"A payment has been made","TransactionID":"d1985d9c-b006-4b38-93bf-4596c7eb9d0c","customerID":3,"merchantID":1,"merchantName":"toko bapak","bankAccountID":1,"amount":5000}
```
Customer has successfully transferred money: :
```
{"level":"info","ts":1698334549.1456852,"caller":"controller/transfer-controller.go:37","msg":"Request transfer money has been initiated","transferID":"9d393793-d00e-44cc-ab83-8e4d68325f92","senderAcountNumber":"12481257","receiverAccountNumber":"128124756","amount":10000}
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

I will use this endpoint /app/transaction/list (middleware protected) to test whether the logout method is successful or not. Take a look at this picture; we have successfully logged out from the system.

While I'm using the customer ID for the logout, I can also use other options such as email or phone number. :)

![alt text](https://i.ibb.co/ZHCTYLF/Screenshot-2023-10-26-190104.jpg)

#### Postman Documentation
https://documenter.getpostman.com/view/29723627/2s9YRFTov9

#### Postman Collection
https://www.postman.com/crimson-crater-616314/workspace/mnc-test/collection/29723627-aed2b8b7-fb73-42b9-be9e-9457d38871ba?action=share&creator=29723627

#### Unit test result
![alt text](https://i.ibb.co/cFxrvLL/unit-test-res.jpg)

#### Database Definition
The SQL scripts for setting up the database are provided in two files: table-definition.sql and dummy-insert.sql.

1.table-definition.sql: This file contains the SQL statements to create the necessary tables for our application. It includes definitions for tables such as merchants, merchantbalances, transactions, and transfer_history.

2.dummy-insert.sql: This file contains SQL statements to insert dummy data into the tables. This could be useful for testing or development purposes.

To set up the database, first run the table-definition.sql script to create the tables, then run the dummy-insert.sql script to populate them with data.

#### Stack
- Gin (Web Framework)
- Paseto (Security concern)
- Zap (Logging)
- Redis (Caching)
