alter table if exists orders
add column payment_id bigint not null,
add constraint fk_pay foreign key(payment_id) references payments(id) on delete cascade;