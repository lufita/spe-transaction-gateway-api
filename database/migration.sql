create table api_clients
(
    id           uuid                     default uuid_generate_v4() not null
        primary key,
    name         varchar(100)                                        not null,
    api_key_hash varchar(255)                                        not null,
    is_active    boolean                  default true,
    rate_limit   integer                  default 60,
    created_at   timestamp with time zone default now(),
    created_by   varchar(64)                                         not null,
    updated_at   timestamp with time zone default now(),
    updated_by   varchar(64)                                         not null,
    client_id    varchar(15)
);

alter table api_clients
    owner to user_sa;

create table auth_logs
(
    id            bigserial
        primary key,
    api_client_id uuid
        references api_clients,
    request_ip    inet,
    endpoint      varchar(255),
    status        varchar(16),
    message       text,
    created_at    timestamp with time zone default now(),
    created_by    varchar(64)
);

alter table auth_logs
    owner to user_sa;

create index idx_auth_logs_client_id
    on auth_logs (api_client_id);

create index idx_auth_logs_ip
    on auth_logs (request_ip);

create table transactions
(
    id                         uuid                        default uuid_generate_v4() not null
        primary key,
    number_billing             varchar(20)                                            not null
        unique,
    request_id                 varchar(32),
    customer_pan               varchar(20),
    amount                     numeric(12, 2),
    transaction_datetime       timestamp(6) with time zone,
    retrieval_reference_number varchar(12),
    customer_name              varchar(100),
    merchant_id                varchar(15),
    merchant_name              varchar(100),
    merchant_city              varchar(100),
    currency_code              varchar(5),
    payment_status             varchar(10),
    payment_description        varchar(100),
    created_at                 timestamp(6) with time zone default now()              not null,
    created_by                 varchar(64)                                            not null
);

alter table transactions
    owner to user_sa;

create table transaction_hist
(
    id             uuid                     default uuid_generate_v4() not null
        constraint transaction_events_pkey
            primary key,
    transaction_id uuid                                                not null
        constraint transaction_events_transaction_id_fkey
            references transactions
            on delete cascade,
    event_type     varchar(64)                                         not null,
    event_data     jsonb,
    created_at     timestamp with time zone default now(),
    created_by     varchar(64)
);

alter table transaction_hist
    owner to user_sa;

create index idx_event_transaction_id
    on transaction_hist (transaction_id);

create table transaction_logs
(
    id             bigserial
        primary key,
    transaction_id uuid not null
        references transactions
            on delete cascade,
    message        text,
    source         varchar(64),
    created_at     timestamp with time zone default now(),
    created_by     varchar(64)
);

alter table transaction_logs
    owner to user_sa;

create index idx_log_transaction_id
    on transaction_logs (transaction_id);

