cd /home/jdumais/Projects/go/src/petServer
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo
docker build -t pet_server:latest .
docker run --mount type=bind,src=/Users/doomer/tmp/pets.json,dst=/Users/doomer/tmp/pets.json -p 8080:8080 pet_server

curl --header "Content-Type: application/json" -X PUT --data '{"pets_collection":{"Buttons":{"age":2,"breed":"Terrier"},"Gracie":{"age":9,"breed":"Spitz"},"Shasta":{"age":9,"breed":"Eskie"}}}' http://localhost:8080/pet
curl http://localhost:8080/pet
curl http://localhost:8080/pet?name=Buttons
curl -X DELETE http://localhost:8080/pet?name=Shasta
curl -X PUT http://localhost:8080/close

docker rm  $(docker ps -q -a)
docker image rm pet_server

