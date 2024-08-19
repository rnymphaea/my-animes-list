# my-animes-list
Basic server written in golang + simple fronted that allows user to create his own list of animes that he watched.

## Summary
- [Technologies](#technologies)
- [Starting the app](#starting-the-app)
- [Allowed endpoints](#allowed-endpoints)
- [To do](#to-do)

## Technologies
- gofiber
- postgreSQL

## Starting the app
- clone the repository
```bash 
$ git clone https://github.com/rnymphaea/my-animes-list.git
```
- start server
```bash
$ cd my-animes-list && go run cmd/main.go 
```

## Allowed endpoints
|      endpoint      |  method  | description                          |
|:------------------:|:--------:|:-------------------------------------|
|        `/`         |  `GET`   | Main page                            |
|      `/login`      |  `GET`   | Login page                           |
|      `/login`      |  `POST`  | Send data to server for log in       |
|     `/signup`      |  `GET`   | Signup page                          |
|     `/signup`      |  `POST`  | Send data to server for sign up      |
|    `/myanimes`     |  `GET`   | Page with user list                  |
|    `/myanimes`     |  `POST`  | Send data to server for add an anime |
|    `/myanimes`     | `DELETE` | Delete anime from list               |
|  `/myanimes/<id>`  |  `GET`   | Get info about anime from list       |
|     `/logout`      |  `GET`   | Log out                              |

## To do
- [ ] database migrations
- [ ] tests
- [ ] update config