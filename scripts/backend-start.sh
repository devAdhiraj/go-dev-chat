echo " ----------- Starting Dev Chat Backend  ------------"

docker-compose up -d kafka postgres --build
sleep 5
docker-compose exec -it kafka kafka-topics --create --topic chat-messages --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
docker-compose up -d chatserver --build
