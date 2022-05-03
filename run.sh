IMAGE_NAME="backend-account-service"
docker build . -f Dockerfile -t $IMAGE_NAME
docker run -i -t --rm --env-file ./.env -p 9527:9527 $IMAGE_NAME
