CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(100) not null,
    created_at timestamp default current_timestamp()
) ENGINE=INNODB;

CREATE TABLE followers(
    user_id int not null,
    foreign key (user_id)
    references users(id)
    on delete cascade,
    follower_id int not null,
    foreign key (follower_id)
    references users(id)
    on delete cascade,
    primary key(user_id, follower_id)
) ENGINE=INNODB;
