echo "=====Building Docker image====="
docker build -t group14/blockchain .
echo "=====Docker image is successfully built====="
echo "=====Running docker-compose file====="
docker-compose up
echo "=====docker-compose is up!====="
