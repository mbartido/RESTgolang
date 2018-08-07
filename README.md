# RESTgolang
A simple RESTful API in golang.
### Install
1. cd cmd/user-list-server
2. go install
3. To run:
  user-list-server --scheme http
### RESTful commands
/user
1. GET:
  curl -i localhost:port/user
2. POST:
  curl -i localhost:port/user -d "{\"name\": \"sampleName\"}" -H "Content-Type: application/json"

/user/{userID}
1. PATCH:
  curl -i localhost:port/user/{id} -X PATCH -H "Content-Type: application/json" -d "{\"name\": \"sampleName\"}"
2. DELETE:
  curl -i localhost:port/user/{id} -X DELETE -H "Content-Type: application/json"
3. GET:
  curl -i localhost:port/user/{id}