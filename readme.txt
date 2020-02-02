cd /home/jdumais/Projects/go/src/ciGatingServer/webServer
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo
docker build -t ci_gating_server:latest .
docker run --mount type=bind,src=/home/jdumais/ciGatingServer/ciGatingServerSettings.json,dst=/dataStore/ciGatingServerSettings.json -p 8080:8080 ci_gating_server

curl --header "Content-Type: application/json" --request PUT --data '{"settings_collection":{"ick":{"enabled":false,"gates_ci_build":false}}}' http://localhost:8080/integrationTest
curl http://localhost:8080/integrationTests?testPath=Gracie
curl http://localhost:8080/integrationTests
curl -X DELETE http://localhost:8080/integrationTests?testPath=Shasta
curl -X PUT http://localhost:8080/close

docker rm  $(docker ps -q -a)
docker image rm ci_gating_server

