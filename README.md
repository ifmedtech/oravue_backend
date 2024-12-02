# oravue_backend

docker build -t oravue .

docker run -e CONFIG_PATH=config/production.yaml oravue

docker build -t oravue .

docker save -o oravue.tar oravue