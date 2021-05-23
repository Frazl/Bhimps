package datamanager

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Setup() *sql.DB {
	database, err := sql.Open("sqlite3", "./bhimps.db")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to initalize the database")
	}
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS userscores (id INTEGER NOT NULL PRIMARY KEY, score INTEGER)")
	statement.Exec()
	statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY, cid INTEGER, score INTEGER)")
	statement.Exec()
	return database
}

// Changes a user's score by a specified amount, can be postiive or negative
func ModifyUserScore(database *sql.DB, user int, changeAmount int) {
	// Ensure user exists in the database
	statement, _ := database.Prepare("SELECT * FROM userscores WHERE id = ?")
	res, _ := statement.Query(user)
	userExists := res.Next()
	res.Close()
	if !userExists {
		statement, err := database.Prepare("INSERT INTO userscores (id, score) VALUES (?, 0)")
		statement.Exec(user)
		if err != nil {
			log.Println(err)
			log.Fatalln("Failed to set user score to inital value")
		}
	}
	statement.Exec(user)
	if changeAmount < 0 {
		decrementUserScore(database, user, changeAmount*-1)
	} else {
		incrementUserScore(database, user, changeAmount)
	}
}

// Increments a user's score by a specified amount
func incrementUserScore(database *sql.DB, user int, amount int) {
	statement, err := database.Prepare("UPDATE userscores SET score = score + ? WHERE id = ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to prepare increment user score query")
	}
	_, err = statement.Exec(amount, user)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to execute increment user score query")
	}
}

// Decrements a user's score by a specified amount
func decrementUserScore(database *sql.DB, user int, amount int) {
	statement, err := database.Prepare("UPDATE userscores SET score = score - ? WHERE id = ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to decrement user score")
	}
	statement.Exec(amount, user)
}

// Get's the score of the requested user
func GetUserScore(database *sql.DB, user int) int {
	statement, err := database.Prepare("SELECT score FROM userscores WHERE id = ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to increment user score")
	}
	res, _ := statement.Query(user)
	res.Next()
	var score int
	res.Scan(&score)
	res.Close()
	return score
}

type UserScore struct {
	ID    int
	Score int
}

func GetUserScores(database *sql.DB, amount int, desc bool) []UserScore {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}
	statement, err := database.Prepare("SELECT id, score FROM userscores WHERE score != 0 ORDER BY score " + order + " LIMIT ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to get user scores")
	}
	res, _ := statement.Query(amount)
	var userscores = make([]UserScore, amount)
	var i int = 0
	for res.Next() {
		var us UserScore
		var id int
		var score int
		res.Scan(&id, &score)
		us.ID = id
		us.Score = score
		userscores[i] = us
		i++
	}
	res.Close()
	return userscores[:i]
}

// Handling message scoreboard

func ModifyMessageScore(database *sql.DB, channel int, message int, changeAmount int) {
	// Ensure user exists in the database
	statement, _ := database.Prepare("SELECT * FROM messages WHERE id = ? AND cid = ?")
	res, err := statement.Query(message, channel)
	if err != nil {
		log.Println("Failed to select * from messages.")
		log.Fatalln(err)
	}
	messageExists := res.Next()
	res.Close()
	if !messageExists {
		statement, err := database.Prepare("INSERT INTO messages (id, cid, score) VALUES (?, ?, 0)")
		statement.Exec(message, channel)
		if err != nil {
			log.Println(err)
			log.Fatalln("Failed to set message score to inital value")
		}
	}
	statement.Exec(message, channel)
	if changeAmount < 0 {
		decrementMessageScore(database, channel, message, changeAmount*-1)
	} else {
		incrementMessageScore(database, channel, message, changeAmount)
	}
}

func incrementMessageScore(database *sql.DB, channel int, message int, amount int) {
	statement, err := database.Prepare("UPDATE messages SET score = score + ? WHERE id = ? AND cid = ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to prepare increment messages score query")
	}
	_, err = statement.Exec(amount, message, channel)
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to execute increment messages score query")
	}
}

func decrementMessageScore(database *sql.DB, channel int, message int, amount int) {
	statement, err := database.Prepare("UPDATE messages SET score = score - ? WHERE id = ? AND cid = ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to decrement messages score")
	}
	statement.Exec(amount, message, channel)
}

type MessageScore struct {
	MessageID string
	ChannelID string
	Score     int
}

func GetMessageScores(database *sql.DB, amount int, desc bool) []MessageScore {
	var order string
	if desc {
		order = "DESC"
	} else {
		order = "ASC"
	}
	statement, err := database.Prepare("SELECT id, cid, score FROM messages WHERE score != 0 ORDER BY score " + order + " LIMIT ?")
	if err != nil {
		log.Println(err)
		log.Fatalf("Failed to get message scores")
	}
	res, _ := statement.Query(amount)
	var messageScores = make([]MessageScore, amount)
	var i int = 0
	for res.Next() {
		var messageScore MessageScore
		var messageID int
		var channelID int
		var score int
		res.Scan(&messageID, &channelID, &score)
		messageScore.MessageID = strconv.Itoa(messageID)
		messageScore.ChannelID = strconv.Itoa(channelID)
		messageScore.Score = score
		messageScores[i] = messageScore
		i++
	}
	return messageScores[:i]
}
