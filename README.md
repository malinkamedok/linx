Build status:
![Status](https://github.com/malinkamedok/linx/actions/workflows/ci.yml/badge.svg)

docker build -t linx:latest .

docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.json