# Тестовое задание Linxdatacenter

Build status:
![Status](https://github.com/malinkamedok/linx/actions/workflows/ci.yaml/badge.svg)

## Задача

Есть 2 файла с данными о продуктах (наименование, цена, рейтинг) в 2-х форматах - CSV и JSON. Примеры в директории [resources](resources).
Необходимо написать программу, которая бы считывала данные из переданного в параметре файла и выводила "самый дорогой продукт" и "с самым высоким рейтингом".
Предусмотреть, что файлы могут быть огромными.


Стек технологий:
- GoLang
- docker
- github

Репозиторий должен содержать Dockerfile для сборки готового приложения в docker среде.
Репозиторий необходимо выложить на github и предоставить ссылку.

## Используемые технологии

- [`Go 1.20`](https://go.dev/)
- Быстрый парсер JSON [`jsonparser`](https://www.github.com/buger/jsonparser)
- mmap библиотека для Go [`mmap-go`](https://github.com/edsrzf/mmap-go)
- Логгер [`logrus`](https://github.com/sirupsen/logrus)
- [`Docker 20.10.21`](https://www.docker.com/) Ссылка на Dockerfile приложения - [Dockerfile](Dockerfile).
  - Использован multi-stage build для разделения процесса сборки

## Запуск приложения

Сборка:
```bash
docker build -t linx:latest .
```

Запуск с парсингом .json файла
```bash
docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.json
```

Запуск с парсингом .csv файла
```bash
docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.csv
```

Запуск с парсингом .json файла и debug логами
```bash
docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.json -v
```

Запуск с парсингом .csv файла и debug логами
```bash
docker run --rm -v $(pwd):$(pwd) -w $(pwd) linx:latest linx --filename ./resources/db.csv -v
```