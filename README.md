# oravue_backend

docker build -t oravue .

docker run -e CONFIG_PATH=config/prod.yaml oravue

docker build -t oravue .

docker save -o oravue.tar oravue

sudo docker build --platform linux/amd64 -t oravue:latest .
