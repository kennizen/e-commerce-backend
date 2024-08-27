create table if not exists customers (
	id bigserial primary key,
	firstname varchar(255) not null,
	middlename varchar(255),
	lastname varchar(255),
	email varchar(1024) unique not null,
	age integer check (age > 0),
	avatar varchar,
	hashed_password varchar,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

create table if not exists products (
	id bigserial primary key,
	title varchar(255) not null,
	description varchar,
	category varchar(255) not null,
	price float not null,
	stock integer not null,
	image varchar not null,
	thumbnail varchar not null,
	rating float default 0.0,
	weight varchar,
	width varchar,
	height varchar,
	depth varchar,
	warranty varchar,
	shipping varchar not null,
	availability varchar not null,
	return_policy varchar not null,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp
);

create table if not exists product_reviews (
	id bigserial primary key,
	review_by bigint not null,
	product_id bigint,
	review varchar,
	rating float not null default 0.0,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(review_by) references customers(id) on delete cascade,
	constraint fk_prod foreign key(product_id) references products(id) on delete set null
);

create table if not exists cart (
	id bigserial primary key,
	cart_of bigint not null,
	items integer default 0,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(cart_of) references customers(id) on delete cascade
);

create table if not exists cart_products_mapping (
	id bigserial primary key,
	cart_id bigint not null,
	product_id bigint not null,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cart foreign key(cart_id) references cart(id) on delete cascade,
	constraint fk_prod foreign key(product_id) references products(id) on delete cascade
);

create table if not exists orders (
	id bigserial primary key,
	order_by bigint not null,
	product_id bigint not null,
	quantity integer default 0,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(order_by) references customers(id) on delete cascade,
	constraint fk_prod foreign key(product_id) references products(id) on delete cascade
);

create table if not exists payments (
	id bigserial primary key,
	payment_by bigint not null,
	payment_method varchar not null,
	payment_status varchar,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(payment_by) references customers(id) on delete cascade
);

create table if not exists favorites (
	id bigserial primary key,
	customer_id bigint not null,
	product_id bigint not null,
	created_at timestamp default current_timestamp,
	updated_at timestamp default current_timestamp,
	constraint fk_cust foreign key(customer_id) references customers(id) on delete cascade,
	constraint fk_prod foreign key(product_id) references products(id) on delete cascade
); 