# Newman tests
Newman tests are used to integration test the API of the application, as well as how it interacts with with the database and backend microservices.

In future, these tests may be expanded to include testing functionality when one or more backend services are down

## Requirements
- npm is required. If you do not have it installed, and are running on Linux or Mac, I recommend installing via nvm: https://github.com/nvm-sh/nvm#troubleshooting-on-linux
- Newman must be installed, installation instructions can be found here: https://github.com/postmanlabs/newman
- All docker containers used by the  project must be running. The docker containers must be fresh, ie there should be no existing entries in the DB.

## How to run
Once newman is installed, you can run these tests with:
```
newman run bank_api_newman_tests.postman_collection.json
```

## Example Results
```
00:01 $ newman run bank_api_newman_tests.postman_collection.json 
newman

bank_api_newman_tests

→ Application - no application ID provided
  GET localhost:8081/api/application?application_id [400 Bad Request, 195B, 34ms]
  ✓  Status code is 400
  ✓  Body matches string

→ Application - application ID does not exist
  GET localhost:8081/api/application?application_id=62cf1e512c8d0c552e1c1ace [404 Not Found, 211B, 6ms]
  ✓  Status code is 404
  ✓  Body matches string

→ Application - application ID format is incorrect
  GET localhost:8081/api/application?application_id=abc [404 Not Found, 190B, 5ms]
  ✓  Status code is 404
  ✓  Body matches string

→ Application with status - No status supplied
  GET localhost:8081/api/applications-with-status [400 Bad Request, 236B, 5ms]
  ✓  Status code is 400
  ✓  Body matches string

→ Applications with status - status is invalid
  GET localhost:8081/api/applications-with-status?status=abc [400 Bad Request, 236B, 4ms]
  ✓  Status code is 400
  ✓  Body matches string

→ Application with status - no applications exist
  GET localhost:8081/api/applications-with-status?status=pending [200 OK, 149B, 5ms]
  ✓  Status code is 200
  ✓  Response matches schema

→ Create Application - No First Name
  POST localhost:8081/api/application [400 Bad Request, 265B, 6ms]
  ✓  Status code is 400
  ✓  Body contains string

→ Create Application - No Last Name
  POST localhost:8081/api/application [400 Bad Request, 263B, 5ms]
  ✓  Status code is 400
  ✓  Body contains string

→ Create Application - Success
  POST localhost:8081/api/application [201 Created, 264B, 7ms]
  ✓  Status code is 201
  ✓  Response matches schema

→ Application with status - Pending exists
  GET localhost:8081/api/applications-with-status?status=pending [200 OK, 339B, 6ms]
  ✓  Status code is 200
  ✓  Response matches schema

┌─────────────────────────┬─────────────────┬─────────────────┐
│                         │        executed │          failed │
├─────────────────────────┼─────────────────┼─────────────────┤
│              iterations │               1 │               0 │
├─────────────────────────┼─────────────────┼─────────────────┤
│                requests │              10 │               0 │
├─────────────────────────┼─────────────────┼─────────────────┤
│            test-scripts │              10 │               0 │
├─────────────────────────┼─────────────────┼─────────────────┤
│      prerequest-scripts │               0 │               0 │
├─────────────────────────┼─────────────────┼─────────────────┤
│              assertions │              20 │               0 │
├─────────────────────────┴─────────────────┴─────────────────┤
│ total run duration: 264ms                                   │
├─────────────────────────────────────────────────────────────┤
│ total data received: 1.05kB (approx)                        │
├─────────────────────────────────────────────────────────────┤
│ average response time: 8ms [min: 4ms, max: 34ms, s.d.: 8ms] │
└─────────────────────────────────────────────────────────────┘

```
