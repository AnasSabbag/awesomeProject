create table if not exists
    users(
        id serial,
        user_name text,
        email text,
        created_at timestamp default now(),
        archived_at timestamp
);
