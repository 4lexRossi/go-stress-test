# go-stress-test

## Construir a Imagem Docker:

### No diretório onde o código Go e o Dockerfile estão localizados, execute o comando:

` docker build -t loadtest . `

## Rodar o Contêiner Docker:

### Para rodar o contêiner com os parâmetros desejados, utilize o seguinte comando:

` docker run loadtest --url=http://google.com --requests=1000 --concurrency=10 `