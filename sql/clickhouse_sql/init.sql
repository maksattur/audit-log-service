create table if not exists event(
    user_id String,
    event_type String,
    timestamp timestamp,
    specific String,
    created_at timestamp  default now()
)
engine = MergeTree
order by (user_id, event_type, timestamp);

insert into event(user_id, event_type, timestamp, specific) VALUES ('104505','bill', now(), '{"customer": "Alex Cross", "items": {"product": "Tea","qty": 6}, "cost": 5862.25}');
insert into event(user_id, event_type, timestamp, specific) VALUES ('28555','bill', now(), '{"customer": "Met Adams", "items": {"product": "Coffee","qty": 8}, "cost": 5841, "discount":568.34}');
insert into event(user_id, event_type, timestamp, specific) VALUES ('7c852e66-3c28-48a3-912a-4e69b4a5973c','system', now(), '{"type": "created", "id": "7c852e66-3c28-48a3-912a-4e69b4a5973c", "email": "example@gmail.com", "other": {"field1": "data1", "field2": "data2", "fieldN":"dataN"}}');
insert into event(user_id, event_type, timestamp, specific) VALUES ('ZyRXUlIu40u0izXJ7EVbHA==','system', now(), '{"type": "deactivated", "id": "7c852e66-3c28-48a3-912a-4e69b4a5973c", "email": "example@gmail.com", "other": {"field1": "data1", "field2": "data2", "fieldN":"dataN"}}');
insert into event(user_id, event_type, timestamp, specific) VALUES ('49440c4a-a1f3-4078-a133-1d3a7353fcf0','system', now(), '{"type": "perform", "client_ip": "80.131.123.137", "request_method": "POST",  "request_payload": {"username":"johnsmith","loginType":"onlineLogin"}, "resource_fragment": "/servlet/login",  "request_uri": "/servlet/login",    "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36", "response_payload": "SUCCESS", "username": "johnsmith"}');