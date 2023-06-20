docker build -t linx:latest .

docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.json