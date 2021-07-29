-- auto-generated definition
create table users
(
    id              bigint unsigned auto_increment,
    phone           varchar(34)                            null,
    email           varchar(89)                            null,
    password        varchar(60)                            null,
    email_verified  int        default 0                   null,
    active          tinyint(1) default 0                   null,
    name            varchar(55)                            null,
    second_name     varchar(55)                            null,
    created_at      timestamp  default (CURRENT_TIMESTAMP) not null,
    updated_at      timestamp  default (CURRENT_TIMESTAMP) null on update CURRENT_TIMESTAMP,
    uuid            varchar(40)                            not null,
    description     varchar(255)                           null,
    nickname        varchar(30)                            null,
    show_subs       tinyint(1)                             null,
    cost            double                                 null,
    notify_email    tinyint(1) default 1                   null,
    notify_telegram tinyint(1)                             null,
    notify_push     int                                    null,
    trial           tinyint(1) default 1                   null,
    language        int                                    null,
    avatar          varchar(144)                           null,
    facebook_id     bigint unsigned                        null,
    google_id       char(21)                               null,
    constraint email_UNIQUE
        unique (email),
    constraint id_UNIQUE
        unique (id),
    constraint phone_UNIQUE
        unique (phone),
    constraint users_facebook_id_uindex
        unique (facebook_id),
    constraint users_google_id_uindex
        unique (google_id),
    constraint users_nickname_uindex
        unique (nickname)
);

alter table users
    add primary key (id);


-- auto-generated definition
create table comments
(
    id           int auto_increment
        primary key,
    user_uuid    varchar(40)                           not null,
    body         varchar(511)                          null,
    type         int                                   not null,
    foreign_uuid varchar(40)                           null,
    created_at   timestamp default (CURRENT_TIMESTAMP) not null,
    updated_at   timestamp default (CURRENT_TIMESTAMP) not null on update CURRENT_TIMESTAMP
);

-- auto-generated definition
create table files
(
    id           bigint unsigned auto_increment
        primary key,
    type         int                                   not null,
    foreign_id   bigint unsigned                       not null,
    dir          varchar(100)                          not null,
    active       int       default 1                   null,
    user_id      bigint                                not null,
    name         varchar(50)                           not null,
    created_at   timestamp default (CURRENT_TIMESTAMP) not null,
    updated_at   timestamp default (now())             not null on update CURRENT_TIMESTAMP,
    uuid         varchar(40)                           not null,
    foreign_uuid varchar(40)                           null,
    user_uuid    varchar(40)                           null,
    constraint files_uuid_uindex
        unique (uuid)
);

-- auto-generated definition
create table friends
(
    id          bigint unsigned auto_increment
        primary key,
    user_uuid   varchar(40)                           not null,
    friend_uuid varchar(40)                           not null,
    type        int                                   null,
    active      tinyint(1)                            null,
    created_at  timestamp default (CURRENT_TIMESTAMP) not null,
    updated_at  timestamp default (CURRENT_TIMESTAMP) not null on update CURRENT_TIMESTAMP
);

-- auto-generated definition
create table hashtags
(
    id         bigint unsigned auto_increment
        primary key,
    tag        varchar(255)                          not null,
    user_id    bigint unsigned                       not null,
    created_at timestamp default (CURRENT_TIMESTAMP) not null,
    constraint hashtags_tag_uindex
        unique (tag)
);

-- auto-generated definition
create table likes
(
    id           bigint unsigned auto_increment
        primary key,
    type         int                                   null,
    foreign_uuid varchar(40)                           null,
    user_uuid    varchar(40)                           null,
    liker_uuid   varchar(40)                           null,
    created_at   timestamp default (CURRENT_TIMESTAMP) not null,
    updated_at   timestamp default (CURRENT_TIMESTAMP) null on update CURRENT_TIMESTAMP,
    active       tinyint(1)                            null,
    constraint likes_type_foreign_id_user_id_uindex
        unique (type, foreign_uuid, user_uuid, liker_uuid)
);

-- auto-generated definition
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

-- auto-generated definition
create table subscribes
(
    id              int auto_increment
        primary key,
    user_uuid       varchar(40)                           not null,
    subscriber_uuid varchar(40)                           not null,
    active          tinyint(1)                            null,
    created_at      timestamp default (CURRENT_TIMESTAMP) null,
    updated_at      timestamp default (CURRENT_TIMESTAMP) null on update CURRENT_TIMESTAMP,
    constraint subscribes_user_id_subscribe_id_uindex
        unique (user_uuid, subscriber_uuid),
    constraint subscribes_user_uuid_subscriber_uuid_uindex
        unique (user_uuid, subscriber_uuid)
);

