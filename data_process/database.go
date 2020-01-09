package data_process

import (
	"database/sql"
	"log"

	"github.com/Mishamba/final_task/model"
	"github.com/golang-migrate/migrate"
	//CAN'T FIND THIS PACKAGES
	//"github.com/golang-migrate/migrate/v4"
	//"github.com/golang-migrate/migrate/v4/database/postgres"
	//_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const connectString = "postgres://tweeter:123@0.0.0.0:5432/tweeter"

func SQLStart(conn *sql.DB) (err error) {
	if conn, err = sql.Open("postgres", connectString); err != nil {
		log.Println("can't connect to db")
		log.Println(err)
		return
	}

	if err = conn.Ping(); err != nil {
		log.Println("can't ping...")
		log.Println(err)
	}

	if m, err := migrate.New("file:///home/go/src/github.com/Mishamba/final_task/db/migrations", "postgres://tweeter:123@0.0.0.0:5432/tweeter_postgres?sslmodedisable"); err == nil {
		err = m.Up()
		log.Println("made migrations successfully")
	} else {
		log.Println("can't make migrations") //always reaches this code
		log.Println(err)
	}

	return
}

func AddUser(user model.User, conn *sql.DB) error {
	info, err := conn.Exec("INSERT INTO Users (username, email, password) VALUE ($1, $2, $3)", user.Name, user.Email, user.Password)
	if n, _ := info.RowsAffected(); err != nil || int(n) != 0 {
		return err
	}

	return nil
}

func PostTweet(tweet string, author int, conn *sql.DB) error {
	info, err := conn.Exec("INSERT INTO Tweets (tweet, authors_id, post_time) VALUE ($1, $2, NOW())", tweet, author)
	if n, _ := info.RowsAffected(); err != nil || int(n) != 0 {
		return err
	}

	return nil
}

func GetTweets(authorID int, conn *sql.DB) ([]model.Tweet, error) {
	var result []model.Tweet
	rows, err := conn.Query("SELECT (tweet, author_ID, upload_time) FROM Tweets")
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var tmp model.Tweet
		err := rows.Scan(&tmp.Tweet, &tmp.AuthorID, &tmp.UploadTime)
		if err != nil {
			return []model.Tweet{}, err
		}

		result = append(result, tmp)
	}

	return result, nil
}

func FindUser(user model.User, conn *sql.DB) (model.User, error) {
	rows, err := conn.Query("SELECT (ID, name, email) FROM Users WHERE ID=$1 OR name=$2 OR email=$3", user.ID, user.Name, user.Email)
	defer rows.Close()
	if err != nil {
		return model.User{}, err
	}

	err = rows.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func GetSubscribers(user model.User, conn *sql.DB) (sub string, err error) {
	rows, err := conn.Query("SELECT subscribers FROM Users WHERE ID=$1 OR name=$2 OR email=$3", user.ID, user.Name, user.Email)
	defer rows.Close()
	if err != nil {
		return sub, err
	}

	if err = rows.Scan(&user.Subscribers); err != nil {
		return sub, err
	}
	return
}

func AddSubscribers(user model.User, newSub int, conn *sql.DB) error {
	rows, err := conn.Query("SELECT subscribers FROM Users WHERE ID=$1", user.ID)
	defer rows.Close()
	if err != nil {
		return err
	}

	if err = rows.Scan(&user.Subscribers); err != nil {
		return err
	}

	user.Subscribers += "*" + string(newSub)

	info, err := conn.Exec("INSERT INTO Users subscribers WHERE ID=$1", user.ID)
	if n, _ := info.RowsAffected(); err != nil || int(n) != 0 {
		return err
	}

	return nil
}
