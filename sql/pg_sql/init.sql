CREATE TABLE IF NOT EXISTS event (
    id SERIAL NOT NULL,
    user_id VARCHAR(255) NOT NULL ,
    event_type VARCHAR(250) NOT NULL,
    timestamp timestamp not null,
    specific jsonb,
    created_at timestamp not null default now(),
    PRIMARY KEY (id)
);

insert into event(user_id, event_type, timestamp, specific) VALUES (gen_random_uuid(),'bill', now(), '{ "customer": "Alex Cross", "items": {"product": "Tea","qty": 6}, "cost": 5862.25}');
insert into event(user_id, event_type, timestamp, specific) VALUES (gen_random_uuid(),'bill', now(), '{ "customer": "Met Adams", "items": {"product": "Coffee","qty": 8}, "cost": 5841, "discount":568.34}');