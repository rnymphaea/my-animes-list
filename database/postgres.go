package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"anime/config"
	"anime/database/models"
)

type storage struct {
	pool *pgxpool.Pool
}

var db storage

func ConnectDB() (*pgxpool.Pool, error) {
	db_url, ok := config.Get("db_url")
	if !ok {
		log.Fatal("db_url env variable not set")
	}
	pool, err := pgxpool.New(context.Background(), db_url)
	if err != nil {
		log.Fatal(err)
	}
	db = storage{pool: pool}
	return pool, err
}

func CreateTables(conn *pgxpool.Pool) {
	createTablesSQL := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(100) UNIQUE NOT NULL,
			password VARCHAR(80) NOT NULL
        );

        CREATE TABLE IF NOT EXISTS anime (
            id SERIAL PRIMARY KEY,
			mal_id INT NOT NULL,
            title VARCHAR(100) NOT NULL,
            description TEXT,
			episodes INT NOT NULL
        );

        CREATE TABLE IF NOT EXISTS user_anime_list (
            id SERIAL PRIMARY KEY,
            user_id INT NOT NULL,
            anime_id INT NOT NULL,
            FOREIGN KEY (user_id) REFERENCES users (id),
            FOREIGN KEY (anime_id) REFERENCES anime (id)
        );
		`

	_, err := conn.Exec(context.Background(), createTablesSQL)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}

func CreateUser(email, password string) error {
	log.Printf("Create user with email: %s, password: %s", email, password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}
	pass := string(hashedPassword)

	_, err = db.pool.Exec(context.Background(), "INSERT INTO users(email, password) values($1, $2)", email, pass)
	if err != nil {
		log.Println("Create user", err)
		return err
	}
	return nil

}

func AddAllAnimes() {
	var data models.Data
	var anime models.Anime
	for i := 0; i < 10000; i++ {
		URL := fmt.Sprintf("https://api.jikan.moe/v4/anime/%d", i)
		response, err := http.Get(URL)
		if response.StatusCode != 200 {
			continue
		}
		if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
			log.Fatal(err)
		} else {
			anime = data.Data
			_, err = db.pool.Exec(context.Background(), "INSERT INTO anime(mal_id, title, description, episodes) VALUES($1, $2, $3, $4)",
				anime.ID, anime.Title, anime.Description, anime.Episodes)
		}
		response.Body.Close()
	}
}

func GetAnimeList(email string) ([]models.Anime, error) {
	var animes []models.Anime
	userID, err := getUserIDbyEmail(email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := db.pool.Query(context.Background(), `
		SELECT anime.mal_id, anime.title, anime.description, anime.episodes 
		FROM user_anime_list
		JOIN anime ON user_anime_list.anime_id = anime.mal_id
		WHERE user_anime_list.user_id = $1`, userID)
	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var anime models.Anime
		if err = rows.Scan(&anime.ID, &anime.Title, &anime.Description, &anime.Episodes); err != nil {
			return nil, err
		}
		animes = append(animes, anime)
	}

	return animes, nil
}

func AddAnime(email, title string) error {
	userID, err := getUserIDbyEmail(email)
	if err != nil {
		log.Println("database/postgres.AddAnime - userID: ", err)
		return err
	}

	animeID, err := getAnimeIDbyTitle(title)
	if err != nil {
		log.Println("database/postgres.AddAnime - animeID: ", err)
		return err
	}

	if animeInList(userID, animeID) {
		log.Println("database/postgres.AddAnime - animeInList: true")
		return errors.New("anime is already in list")
	}

	_, err = db.pool.Exec(context.Background(), "INSERT INTO user_anime_list(user_id, anime_id) VALUES($1, $2)", userID, animeID)
	if err != nil {
		log.Println("database/postgres.AddAnime - exec: ", err)
		return err
	} else {
		return nil
	}
}

func animeInList(userID, animeID int) bool {
	var id int
	row := db.pool.QueryRow(context.Background(), "SELECT id FROM user_anime_list WHERE user_id = $1 AND anime_id = $2", userID, animeID)

	if err := row.Scan(&id); err != nil {
		log.Println("id = ", id)
		log.Println("animeInList: ", animeID, userID, err)
		if errors.Is(err, pgx.ErrNoRows) {
			return false
		} else {
			return true
		}
	}
	return true
}

func getAnimeIDbyTitle(title string) (int, error) {
	var id int
	row := db.pool.QueryRow(context.Background(), "SELECT mal_id FROM anime WHERE title = $1", title)
	if err := row.Scan(&id); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func GetAnimebyID(id int) (models.Anime, error) {
	var anime models.Anime
	row := db.pool.QueryRow(context.Background(), "SELECT mal_id, title, description, episodes FROM anime WHERE mal_id=$1", id)
	if err := row.Scan(&anime.ID, &anime.Title, &anime.Description, &anime.Episodes); err != nil {
		log.Println("database/postgres.GetAnimebyID: ", err)
		return anime, err
	} else {
		return anime, nil
	}
}

func getUserIDbyEmail(email string) (int, error) {
	var id int
	row := db.pool.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func CheckPassword(email, password string) (bool, error) {
	var passInDB string
	row := db.pool.QueryRow(context.Background(), "SELECT password FROM users WHERE email = $1", email)
	if err := row.Scan(&passInDB); err != nil {
		log.Println("database/postgres.CheckPassword: ", err)
		return false, err
	}
	check := bcrypt.CompareHashAndPassword([]byte(passInDB), []byte(password))
	return check == nil, nil
}
