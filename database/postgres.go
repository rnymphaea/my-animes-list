package database

import (
	"anime/database/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type storage struct {
	pool *pgxpool.Pool
}

var db storage

func ConnectDB() (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:psql@localhost:5432/animes_list")
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return err
	}
	pass := string(hashedPassword)

	_, err = db.pool.Exec(context.Background(), "INSERT INTO users(email, password) values($1, $2)", email, pass)
	if err != nil {
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
		log.Println(err)
	}
	animeID, err := getAnimeIDbyTitle(title)
	if err != nil {
		log.Println(err)
	}
	_, err = db.pool.Exec(context.Background(), "INSERT INTO user_anime_list(user_id, anime_id) VALUES($1, $2)", userID, animeID)
	if err != nil {
		return err
	} else {
		log.Println("Success")
		return nil
	}
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

func getUserIDbyEmail(email string) (int, error) {
	var id int
	row := db.pool.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	} else {
		return id, nil
	}
}
