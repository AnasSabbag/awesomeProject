begin;

create table user_role
(
    id          serial
        primary key,
    name        text,
    description text,
    created_at  timestamp with time zone default now(),
    archived_at timestamp with time zone,
    updated_at  timestamp with time zone
);
create table users
(
    id             serial
        primary key,
    name           text                                   not null,
    email          text                                   not null,
    username       text                                   not null,
    password       text                                   not null,
    created_on     timestamp with time zone default now() not null,
    last_login     timestamp with time zone default now() not null,
    is_admin       boolean                                not null,
    is_deactivated boolean                  default false,

    role_id        integer                  default 1
        references user_role,
    archived_at    timestamp with time zone
);

create table user_permission
(
    id          serial
        primary key,
    name        text,
    description text,
    is_deleted  boolean default false,
    created_at timestamp with time zone default now()
);
create table user_role_permission_relation
(
    id            serial
        primary key,
    role_id       integer
        references user_role,
    permission_id integer
        references user_permission
);


create table sessions
(
    id          text not null
        primary key,
    user_id     integer
        references users,
    expiry_time timestamp with time zone default (now() + '00:15:00'::interval),
    archived_at timestamp with time zone,
    created_at  timestamp with time zone default now()
);


commit;