drop table addresses;

alter table orders
drop column address_used,
drop constraint fk_add;