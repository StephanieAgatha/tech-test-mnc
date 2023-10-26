--bank accounts
create table bankaccounts
(
    id             serial
        primary key,
    bank_id        integer
        references banks,
    customer_id    integer
        references customers,
    account_number varchar(255)
        unique,
    balance        integer default 0,
    unique (customer_id, bank_id)
);

--table banks
create table banks
(
    id   serial
        primary key,
    name varchar(255)
        unique
);

--table customers
create table customers
(
    id           serial
        primary key,
    name         varchar(255) not null,
    email        varchar(255) not null
        unique,
    password     varchar(255) not null,
    created_at   timestamp default CURRENT_TIMESTAMP,
    updated_at   timestamp default CURRENT_TIMESTAMP,
    phone_number varchar(14)
        constraint unique_phone_number
            unique
);

--table merchant balances
create table merchantbalances
(
    id          serial
        primary key,
    merchant_id integer
        references merchants,
    balance     integer default 0 not null
);

--table merchants
create table merchants
(
    id         serial
        primary key,
    name       varchar(255) not null,
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);

--table transactions
create table transactions
(
    id              serial
        primary key,
    customer_id     integer
        references customers,
    merchant_id     integer
        references merchants,
    bank_account_id integer not null
        references bankaccounts,
    amount          integer not null,
    created_at      timestamp default CURRENT_TIMESTAMP,
    updated_at      timestamp default CURRENT_TIMESTAMP
);

--table transfer history
create table transfer_history
(
    id                      serial
        primary key,
    sender_account_number   varchar(15)
        references bankaccounts (account_number),
    receiver_account_number varchar(15)
        references bankaccounts (account_number),
    amount                  integer,
    transfer_timestamp      timestamp default CURRENT_TIMESTAMP,
    transaction_id          varchar(255)
);
