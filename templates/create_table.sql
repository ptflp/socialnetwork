create table posts
(
    id          bigint unsigned auto_increment
        primary key,
    body        varchar(1024)                         not null,
    file_id     int                                   null,
    active      int       default 0                   not null,
    created_at  timestamp default (CURRENT_TIMESTAMP) not null,
    updated_at  timestamp default (CURRENT_TIMESTAMP) null on update CURRENT_TIMESTAMP,
    user_id     bigint                                not null,
    type        int       default 1                   not null,
    is_reposted int       default 0                   not null,
    repost_id   int                                   null,
    uuid        varchar(40)                           not null,
    file_uuid   varchar(40)                           not null,
    user_uuid   varchar(40)                           not null,
    price       int                                   null
);