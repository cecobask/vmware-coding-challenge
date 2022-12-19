# vmware-coding-challenge
Write a Golang based HTTP server which accepts GET requests with input parameter as "sortKey"
and "limit". The server queries three URLs mentioned below, combines the results from all three
URLs, sorts them by the sortKey and returns the response. The server should also limit the number
of items in the API response to input parameter "limit".

## URLs
* https://raw.githubusercontent.com/assignment132/assignment/main/duckduckgo.json
* https://raw.githubusercontent.com/assignment132/assignment/main/google.json
* https://raw.githubusercontent.com/assignment132/assignment/main/wikipedia.json

## Parameters
| Parameter | Type    | Example                       |
|-----------|---------|-------------------------------|
| sortKey   | String  | relevanceScore / views        |
| limit     | Integer | Greater than 1, less than 200 |

## Response Format
```json
{
  "data": [
    {
      "url": "www.yahoo.com/abc6",
      "views": 6000,
      "relevanceScore": 0.6
    },
    ...
  ],
  "count": 6
}
```

## Requirements
- [x] The server should query URLs concurrently
- [x] Server should have re-try mechanism and error handling on failures while querying URLs
- [ ] Code should have unit tests
- [x] Provide README for testing and deployment
- [x] Provide deployment manifests to deploy the service in a Kubernetes cluster by following best
  practices

## Run
### Kubernetes
```shell
make kubernetes-deploy
```
### Source code
```shell
make run
```
### Docker
```shell
make docker-build
make docker-run
```

## Usage
```shell
curl -sX GET 'localhost:8080/pagedata?sortKey=relevanceScore&limit=15' | json_pp
```