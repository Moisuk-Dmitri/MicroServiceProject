# MicroServiceProject

Учебный проект для практики микросервисной архитектуры на Go.

Проект представляет собой монорепозиторий с несколькими сервисами, которые взаимодействуют между собой через:

* REST (HTTP)
* gRPC
* Kafka
* PostgreSQL

---

# Архитектура проекта

Система состоит из 3 сервисов:

## Auth Service

Сервис авторизации и регистрации.

### Возможности

* регистрация пользователей
* login пользователей
* генерация токена
* gRPC проверка токена
* отправка события `user.created` в Kafka

### Использует

* HTTP
* gRPC
* PostgreSQL
* Kafka Producer

---

## Blog Service

Сервис статей.

### Возможности

* создание постов
* получение постов
* middleware авторизации
* проверка токена через Auth Service по gRPC

### Использует

* HTTP
* gRPC Client
* PostgreSQL

---

## Notification Service

Сервис нотификаций.

### Возможности

* прослушивание Kafka topic
* обработка событий создания пользователя

### Использует

* Kafka Consumer

---

# Схема взаимодействия

## Авторизация

```text
Client
  ↓ HTTP
Auth Service
  ↓ PostgreSQL
Database
```

## Проверка токена

```text
Client
  ↓ HTTP
Blog Service
  ↓ gRPC
Auth Service
```

## События Kafka

```text
Auth Service
  ↓ Kafka Producer
Kafka
  ↓ Kafka Consumer
Notification Service
```

---

# Стек технологий

* Go
* Chi Router
* gRPC
* Apache Kafka
* PostgreSQL
* Docker
* Docker Compose

---

# Структура проекта

```text
MicroServiceProject/
├── cmd/
│   ├── auth/
│   ├── blog/
│   └── notification/
│
├── internal/
│   ├── auth/
│   ├── blog/
│   ├── notification/
│   └── platform/
│
├── migrations/
├── proto/
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

---

# Миграции

Проект использует SQL миграции для создания таблиц PostgreSQL.

Миграции находятся в директории:

```text
migrations/
```

---

# Makefile

Проект содержит Makefile для упрощения запуска сервисов и Docker инфраструктуры.

## Запуск сервисов

Запуск Auth Service:

```bash
make auth
```

Запуск Blog Service:

```bash
make blog
```

Запуск Notification Service:

```bash
make notification
```

---

## Работа с Docker Compose

Запуск только Kafka:

```bash
make kafka-up
```

Запуск только PostgreSQL:

```bash
make postgres-up
```

Запуск Docker Compose в фоне:

```bash
make docker-up
```

Остановка Docker Compose:

```bash
make docker-down
```

Полный запуск проекта с пересборкой контейнеров:

```bash
make project-up
```

Полная остановка проекта:

```bash
make project-down
```

Makefile позволяет быстро запускать отдельные сервисы и инфраструктуру без необходимости вручную вводить длинные команды.

---

# Основные особенности проекта

* монорепозиторий
* микросервисная архитектура
* REST API
* gRPC communication
* event-driven communication через Kafka
* middleware авторизации
* Docker Compose инфраструктура
* разделение сервисов и слоев
* dependency injection
* graceful shutdown

---

# Цель проекта

Практика backend-разработки и микросервисной архитектуры на Go с использованием:

* HTTP
* gRPC
* Kafka
* PostgreSQL
* Docker
