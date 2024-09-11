create table if not exists addresses (
	id bigserial primary key,
	country varchar not null,
	state varchar not null,
	address varchar not null,
	zipcode varchar not null,
	phone_number varchar not null unique,
	address_of bigint not null,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(address_of) references customers(id) on delete cascade	
);

alter table if exists orders
add column address_used bigint,
add constraint fk_add foreign key(address_used) references addresses(id) on delete set null;