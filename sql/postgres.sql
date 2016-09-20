drop table if exists sensorlogs;
drop table if exists sensors;
drop table if exists users cascade;

create table sensors (
	id serial unique,
	created timestamp with time zone default now(),
	name text unique not null,
	type text not null
);

create table sensorlogs (
	id serial unique,
	created timestamp with time zone default now(),
	sensorid int references sensors (id) on delete cascade,
	metrics hstore
);

create table users (
        id serial unique,
        created timestamp with time zone default now(),
        fname text not null,
        lname text not null,
        email text not null,
        hash text not null,
        username text unique not null,
        admin bool default false not null
);

insert into sensors (name, type) values ('GreenHouse', 'Weather') returning id;
insert into sensors (name, type) values ('House', 'Weather') returning id;

create or replace function hash(pass text) returns text as $$
        select crypt(pass, gen_salt('bf', 10));
$$ language sql;

insert into users (fname, lname, username, hash, email, admin) values ('Aaron', 'Bieber', 'aaron', hash('omgSnakes'), 'aaron@bolddaemon.com', true);


