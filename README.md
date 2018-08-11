# RESTgolang
A simple RESTful API in golang.
### Install
1. cd swagger3/cmd/user-list-server
2. go install
3. To run:
  user-list-server --scheme http /port 8000
### RESTful commands
/user
1. GET:
  curl -i localhost:8000/user
2. POST:
  curl -i localhost:8000/user -d "{\"name\": \"sampleName\"}" -H "Content-Type: application/json"

/user/{userID}
1. PATCH:
  curl -i localhost:8000/user/{id} -X PATCH -H "Content-Type: application/json" -d "{\"name\": \"sampleName\"}"
2. DELETE:
  curl -i localhost:8000/user/{id} -X DELETE -H "Content-Type: application/json"
3. GET:
  curl -i localhost:8000/user/{id}

To see changes go to swagger3/cmd/user-list-server/users.json.