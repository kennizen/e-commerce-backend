alter table if exists payments
add column amount numeric(10,4) not null default 0;