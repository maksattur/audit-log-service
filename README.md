# Audit log service

Greetings,

I'm pleased to introduce the Audit Log Service, the system designed to process event messages from a variety of sources. This README will provide you with a technical overview of the system's functionality and design, which you can use as reference on GitHub.

## Functional requirements
Our audit service must meet the following requirements:

- **Receiving event messages:** Our system needs to receive event messages from other services.
- **HTTP endpoint:** Our system must provide an HTTP endpoint. With filtering of the requested data by certain fields.
- **Event message format:** Event messages are expected to be in JSON format, consisting of a set of required common invariant fields and specific variant fields.
- **Authentication:** Our system must provide authentication to the HTTP endpoint

## Preliminary considerations
Before diving into the details of system design, it is important to consider two critical aspects:

- **Data providers at the transport layer:** Given the significant influx of data, I'm guessing it could be a message broker. Apache Kafka, NATS, RabbitMQ and similar technologies are well suited for this role.
- **Selecting a database:** For intensive data recording, OLAP databases such as Amazon Redshift, Google BigQuery, ClickHouse, etc. are well suited.

## Technology stack
Given the above considerations and the lack of special requirements, I settled on the following technology stack for our implementation:

- **Message Broker:** Apache Kafka was chosen as our message broker due to its reliability and scalability.
- **Database:** ClickHouse is capable of processing large volumes of data and supporting efficient analytical processing.

I've implemented our code following the principles of a clean architecture. The Adapter design pattern is cleverly used to connect various components of the system. This design choice ensures flexibility, making it a straightforward process to integrate alternative message brokers or databases down the road. All that's needed is the implementation of the contracts outlined in code. For example, the code implements postgres database integration (although it is not entirely suitable for such a role).

## System design
**As a result, the design of our simple system will look like this:**

![System design](https://github.com/maksattur/audit-log-service/blob/main/design.png)


## Event structure
I defined the transport model for the Event structure as follows:
```
type Event struct {
    Common   CommonFields    `json:"common"`
    Specific json.RawMessage `json:"specific"`
}

type CommonFields struct {
    UserID    string    `json:"user_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
}
```
As you can see from this structure, I've defined the following common invariant fields:

- **user_id**
- **event_type**
- **timestamp**

## API design
Documentation detailing the design of API can be found in the OpenAPI specification (swagger.yml). Additionally, I'd like to highlight the introduction of a query parameter called limit. By default, its value is set to 10,000 records. This parameter has been introduced with the purpose of safeguarding the service against overload and potential DoS attacks.

## Authentication
Authentication in our system is implemented using JWT. To obtain a token, you must authenticate with the single default user. The credentials are as follows:

- **Username:** "bambino"
- **Password:** "qwerty123"

## Curl
Authentication
```
curl -i --location 'localhost:8080/auth' \
--header 'Content-Type: application/json' \
--data '{
    "login": "bambino",
    "password": "qwerty123"
}'
```

Get events
```
curl -i --location 'localhost:8080/events?user_id=104505&event_type=bill&from=2023-09-23T14%3A38%3A25Z&to=2023-09-30T14%3A43%3A13Z&limit=1000' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTU2NjEzODF9.FhPq5aoDT2VCeCv2AOWHHahrULFj9jEA3V9eB1PwV2U'
```

## Running the application
Commands for running the application are provided in the Makefile:

- **make dc:** Launches the infrastructure and the application in Docker containers. **It includes a check for the availability of the infrastructure and records a set of default events in the database**
- **make run:** Starts the application on your local machine. You will need to configure all the necessary data as described in the configuration file.
- **make test** Runs the tests for the application.

## TODO
- It would be necessary to cover the code with a large number of tests to ensure comprehensive testing of the application functionality. Now the handler is covered with tests.
- It would be necessary to implement more detailed logging throughout the application to make troubleshooting and monitoring easier.
- It would be good to implement JWT authentication with refresh tokens to improve security and improve user experience. Alternatively, we can use third party services such as **[Keycloak](https://www.keycloak.org)** for advanced authentication and authorization capabilities.
- It would be good to optimize data recording in ClickHouse by implementing batch processing. This can improve the efficiency of writing large amounts of data.