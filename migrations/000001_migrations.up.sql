create type user_type as enum ('vendor', 'receiver');

create table users 
(
        id uuid primary key default gen_random_uuid(),
        user_name varchar not null,
        full_name varchar not null,
        user_type varchar not null,
        bio text,
        password var2char not null,
        email varchar not null,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
        deleted_at TIMESTAMP WITH TIME ZONE
);


create table refresh_token 
(
    id UUID PRIMARY KEY default gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    token TEXT UNIQUE,
    revoked BOOLEAN,
    created_at TIMESTAMP default current_timestamp,
    updated_at TIMESTAMP default current_timestamp,
    deleted_at TIMESTAMP
)