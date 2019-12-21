package data_process

import (
	"database/sql"

	"github.com/Mishamba/FinalTaks/model"
	_ "github.com/lib/pq"
)

const connectString = "host=127.0.0.1 port=5432 user=tweeter dbname=tweet_postgres sslmode=disable"

func ConnectToSQL(conn *sql.DB) (err error) {
	conn, err = sql.Open("postgres", connectString)
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
