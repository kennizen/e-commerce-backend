alter table if exists addresses
add constraint addresses_phone_number_key unique (phone_number);