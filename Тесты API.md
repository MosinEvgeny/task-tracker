# Task Tracker API - Тест-кейсы для Postman

## Общие замечания

* URL: <http://localhost:8080> (или ваш настроенный адрес)
* Content-Type: application/json (для всех запросов с телом)
* Authorization: Bearer \<token> (для защищенных маршрутов) - токен, полученный после успешного логина

## 1. Пользователи

### 1.1 Регистрация пользователя (POST /register)

Запрос:

```json
{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password"
}
```

Ожидаемый ответ:

* Код: 201 Created
* JSON: (Объект пользователя с сгенерированным ID и захешированным паролем)

```json
{
    "id": "008e43-047-4e3-96b-e27214235",
    "username": "testuser",
    "email": "test@example.com",
    "password": "$2a$cu560WIF.bqUDvfOLbc/PPiD/mgrERLWcC9."
}
```

Негативные тесты:

* Неверный формат email (код 400 Bad Request, сообщение об ошибке)
* Email уже существует (код 400 Bad Request, сообщение об ошибке)
* Отсутствуют обязательные поля (код 400 Bad Request)

### 1.2 Вход пользователя (POST /login)

Запрос:

```json
{
    "email": "test@example.com",
    "password": "password"
}
```

Ожидаемый ответ:

* Код: 200 OK
* JSON: (Токен)

```json
{
    "token": "yJleHAiOjE3MzkzMDY3MTQsInVzZXJfaWQiOiIwODA4MWU0My0wMTQ3LTQwZTMtOWQ2Yi1lMjcy",
    "refresh_token": "740efa-1b2e-4668-9b0a-c531d"
}
```

Негативные тесты:

* Неверный email или пароль (код 401 Unauthorized)
* Отсутствуют обязательные поля (код 400 Bad Request)

### 1.3 Получение пользователя по ID (GET /users/{id})

Запрос: (Необходимо добавить заголовок Authorization)
Ожидаемый ответ:

* Код: 200 OK
* JSON: (Объект пользователя)

```json
{
    "id": "...",
    "username": "testuser",
    "email": "<test@example.com>",
    "password": "..."
}
```

Негативные тесты:

* Неверный ID (код 400 Bad Request)
* Пользователь не найден (код 404 Not Found)
* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)

### 1.4 Обновление пользователя (PUT /users/{id})

Запрос: (Необходимо добавить заголовок Authorization)

```json
{
    "username": "newuser",
    "email": "<new@example.com>"
}
```

Ожидаемый ответ:

* Код: 200 OK
* JSON: (Обновленный объект пользователя)

```json
{
    "id": "...",
    "username": "newuser",
    "email": "<new@example.com>",
    "password": "..."
}
```

Негативные тесты:

* Неверный ID (код 400 Bad Request)
* Пользователь не найден (код 404 Not Found)
* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)
* Неверный формат запроса (код 400 Bad Request)

### 1.5 Удаление пользователя (DELETE /users/{id})

Запрос: (Необходимо добавить заголовок Authorization)
Ожидаемый ответ:

* Код: 204 No Content

Негативные тесты:

* Неверный ID (код 400 Bad Request)
* Пользователь не найден (код 404 Not Found)
Отсутствует заголовок Authorization (код 401 Unauthorized)
Неверный токен (код 401 Unauthorized)

### 1.6 Обновление токена (POST /refresh)

Запрос:

```json
{
    "refresh_token": "740efa-1b2e-4668-9b0a-c531d" // (refresh_token, полученный после логина)
}
```

Ожидаемый ответ:

* Код: 200 OK
* JSON: (Новый токен)

```json
{
    "token": "eyJleHAiOjE3MzkzMDcxMjEsInVzZXJfaWQiOiIwODA4MWU0My0wMTQ3LTQwZTMtOWQ2Y"
}
```

Негативные тесты:

* Неверный refresh token (код 401 Unauthorized)
* Срок действия refresh token истек (код 401 Unauthorized)
* Отсутствует refresh token (код 400 Bad Request)

### 1.7 Отзыв всех refresh токенов (POST /users/revoke)

Запрос: (Необходимо добавить заголовок Authorization)

Ожидаемый ответ:

* Код: 204 No Content

## 2. Задачи

### 2.1 Создание задачи (POST /tasks)

Запрос: (Необходимо добавить заголовок Authorization)

```json
{
    "title": "New Task",
    "description": "Task Description",
    "due_date": "2024-03-15T12:00:00Z",
    "user_id": "..." // (ID пользователя)
}
```

Ожидаемый ответ:

* Код: 201 Created
* JSON: (Объект задачи)

Негативные тесты:

* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)
* Неверный формат запроса (код 400 Bad Request)
* Отсутствует обязательное поле (код 400 Bad Request)

### 2.2 Получение задачи (GET /tasks/{id})

(Аналогично пункту 1.3, замените “пользователя” на “задачу”)

### 2.3 Обновление задачи (PUT /tasks/{id})

Запрос: (Необходимо добавить заголовок Authorization)

```json
{
    "title": "Updated Task",
    "description": "Updated Description",
    "due_date": "2024-03-16T12:00:00Z"
}
```

Ожидаемый ответ:

* Код: 200 OK
* JSON: (Объект задачи)
* Негативные тесты:
* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)
* Неверный формат запроса (код 400 Bad Request)

### 2.4 Удаление задачи (DELETE /tasks/{id})

(Аналогично пункту 1.5, замените “пользователя” на “задачу”)

## 3. Метки

### 3.1 Создание метки (POST /labels)

Запрос: (Необходимо добавить заголовок Authorization)

```json
{
    "name": "Important",
    "color": "#FF0000",
    "user_id": "..." // (ID пользователя)
}
```

Ожидаемый ответ:

* Код: 201 Created
* JSON: (Объект метки)

Негативные тесты:

* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)
* Неверный формат запроса (код 400 Bad Request)
* Отсутствует обязательное поле (код 400 Bad Request)
* Неверный формат цвета (код 400 Bad Request)

### 3.2 Получение метки (GET /labels/{id})

(Аналогично пункту 1.3, замените “пользователя” на “метку”)

### 3.3 Обновление метки (PUT /labels/{id})

Запрос: (Необходимо добавить заголовок Authorization)

```json
{
    "name": "Very Important",
    "color": "#00FF00"
}
```

Ожидаемый ответ:

* Код: 200 OK
* JSON: (Объект метки)

Негативные тесты:

* Отсутствует заголовок Authorization (код 401 Unauthorized)
* Неверный токен (код 401 Unauthorized)
* Неверный формат запроса (код 400 Bad Request)
* Неверный формат цвета (код 400 Bad Request)

### 3.4 Удаление метки (DELETE /labels/{id})

(Аналогично пункту 1.5, замените “пользователя” на “метку”)

## Примечания

Замените ... на фактические значения.
Обратите внимание на форматы дат и времени (ISO 8601).
Проверьте все негативные сценарии для каждого запроса.
