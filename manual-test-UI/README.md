# Tic-Tac-Toe manual test UI

Готовый браузерный UI для ручного тестирования REST API из учебного Go-проекта.

## Почему здесь есть `server.py`

Go API слушает `localhost:8080` и не включает CORS. UI поэтому запускается отдельным маленьким Python proxy на `localhost:3000`: браузер ходит на тот же origin (`/api/...`), а proxy пересылает запросы в Go backend.

Требования: Python 3.9+ (только стандартная библиотека).

## Запуск

1. Запусти свой Go backend обычным способом. Он должен быть доступен на `http://127.0.0.1:8080`.
2. В этой папке выполни:

```bash
python3 server.py
```

3. Открой:

```text
http://127.0.0.1:3000
```

Если backend использует другой host/port:

```bash
TICTACTOE_BACKEND_HOST=127.0.0.1 TICTACTOE_BACKEND_PORT=8080 python3 server.py
```

## Возможности

UI покрывает все HTTP-маршруты, зарегистрированные текущим `web/router/router.go`:

- signup;
- login;
- access token refresh;
- refresh token rotation;
- `GET /user/me`;
- `GET /user/{id}` — поиск пользователя по UUID;
- создание PvE и PvP игр;
- `GET /game/{id}`;
- join PvP;
- клики по доске → реальные `POST /game/{id}`;
- `GET /game/available`;
- `GET /game/history`;
- `GET /leaderboard` с выводом лидерборда в отдельной таблице;
- полный HTTP request/response log, включая статус и latency;
- копирование Game ID;
- автоматический polling PvP;
- независимые `sessionStorage`-сессии в разных вкладках.

Клик по строке лидерборда автоматически подставляет UUID пользователя в User lookup и вызывает `GET /user/{id}`.

## PvP

Открой две вкладки `http://127.0.0.1:3000/`.

В каждой вкладке залогинься под отдельным пользователем. Благодаря `sessionStorage` токены вкладок независимы.

В первой вкладке создай PvP игру. Во второй нажми `Available`, выбери игру и `JOIN`. После этого можно играть из двух вкладок.

## HTTP log

Каждое действие через API показывается в журнале справа: метод, endpoint, HTTP status, latency, тело запроса и полный ответ backend.

## Авторизация

UI специально не хранит access/refresh tokens в `localStorage`: каждая вкладка имеет свою `sessionStorage`-сессию. Это удобно именно для ручной проверки PvP двумя разными пользователями.

`Logout` в UI очищает локальную сессию браузера; отдельного logout endpoint в текущем Go API нет.
