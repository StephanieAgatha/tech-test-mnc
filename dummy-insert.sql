--table bank accounts
INSERT INTO bankaccounts (bank_id, customer_id, account_number, balance)
VALUES (1, 1, '12481257', 415000);

INSERT INTO bankaccounts (bank_id, customer_id, account_number, balance)
VALUES (2, 2, '128124756', 70000);

INSERT INTO bankaccounts (bank_id, customer_id, account_number, balance)
VALUES (3, 1, '12371246', 44000);


--table banks
INSERT INTO banks (name)
VALUES ('BCA');

INSERT INTO banks (name)
VALUES ('BRI');

INSERT INTO banks (name)
VALUES ('BNI');

--table customers
INSERT INTO customers (id, name, email, password, phone_number)
VALUES (3, 'sora', 'sorsoraaa@gmail.com', '$2a$10$8nPv0q/bQdmNhnv/q5l23esz8/brfjAXAmJ33e4Bpb6PEZxlIPly.', '085156810932');

INSERT INTO customers (id, name, email, password, phone_number)
VALUES (5, 'awd', 'awd@gmail.com', '$2a$10$9jodaU.vQ4oDVRD5H5rGB.0D1r2FUNgnFUeLjDknnW5.r5DMBenHu', '085156810912');

INSERT INTO customers (id, name, email, password,phone_number)
VALUES (13, 'test3', 'test3@gmail.com', '$2a$10$136klYVNU5cmGNG8mCtioe3k3PGymdTZ3FkHrHQ6F1JzKXRXg9nLO', '0851134681092');

--table merchant balances
INSERT INTO merchantbalances (id, merchant_id, balance)
VALUES (1, 1, 137000);

--table merchants
INSERT INTO merchants (id, name)
VALUES (1, 'toko bapak');

INSERT INTO merchants (id, name)
VALUES (2, 'haji barokah');

--table transactions
INSERT INTO transactions (id, customer_id, merchant_id, bank_account_id, amount, created_at, updated_at)
VALUES (2, 3, 1, 1, 15000, '2023-10-25 16:55:31.751896', '2023-10-25 16:55:31.751896');

INSERT INTO transactions (id, customer_id, merchant_id, bank_account_id, amount, created_at, updated_at)
VALUES (3, 3, 1, 1, 1000, '2023-10-26 00:11:11.062746', '2023-10-26 00:11:11.062746');

INSERT INTO transactions (id, customer_id, merchant_id, bank_account_id, amount, created_at, updated_at)
VALUES (4, 3, 1, 1, 5000, '2023-10-26 00:13:49.048078', '2023-10-26 00:13:49.048078');

INSERT INTO transactions (id, customer_id, merchant_id, bank_account_id, amount, created_at, updated_at)
VALUES (9, 3, 1, 1, 5000, '2023-10-26 01:33:29.228417', '2023-10-26 01:33:29.228417');


--table transactions history
INSERT INTO transfer_history (id, sender_account_number, receiver_account_number, amount, transfer_timestamp, transaction_id)
VALUES (15, '12481257', '12371246', 10000, '2023-10-26 10:49:51.391191', '466d6803-1fb4-4cca-a630-bf1322c36bb0');

INSERT INTO transfer_history (id, sender_account_number, receiver_account_number, amount, transfer_timestamp, transaction_id)
VALUES (16, '12481257', '128124756', 10000, '2023-10-26 10:52:39.333791', 'e3abcc75-4fb2-4172-b4ae-5a8423fdc477');

INSERT INTO transfer_history (id, sender_account_number, receiver_account_number, amount, transfer_timestamp, transaction_id)
VALUES (17, '12481257', '128124756', 10000, '2023-10-26 11:15:20.636468', 'e4250937-912e-4f51-ac50-363e790c02fd');
