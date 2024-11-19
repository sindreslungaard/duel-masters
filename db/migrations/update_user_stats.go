package migrations

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateUserStats(db *mongo.Database) {
    log.Println("Starting UpdateUserStats migration")

    ctx := context.Background()

    userStats := make(map[string]*UserStat)

    cur, err := db.Collection("duels").Find(ctx, bson.M{})
    if err != nil {
        log.Fatal("Failed to fetch duels:", err)
    }
    defer cur.Close(ctx)

    for cur.Next(ctx) {
        var duel Duel
        if err := cur.Decode(&duel); err != nil {
            log.Println("Failed to decode duel:", err)
            continue
        }

        processPlayer(duel.P1, duel.Winner, userStats)
        processPlayer(duel.P2, duel.Winner, userStats)
    }

    if err := cur.Err(); err != nil {
        log.Fatal("Cursor error:", err)
    }

    for uid, stats := range userStats {
        filter := bson.D{{Key: "uid", Value: uid}}
        update := bson.D{
            {Key: "$set", Value: bson.D{
                {Key: "total_games_played", Value: stats.TotalGamesPlayed},
                {Key: "games_won", Value: stats.GamesWon},
                {Key: "games_lost", Value: stats.GamesLost},
                {Key: "win_rate", Value: stats.WinRate},
            }},
        }

        _, err := db.Collection("users").UpdateOne(ctx, filter, update)
        if err != nil {
            log.Println("Failed to update user", uid, ":", err)
        } else {
            log.Println("Updated user stats for", uid)
        }
    }

    log.Println("UpdateUserStats migration completed")
}


type Duel struct {
	UID    string `bson:"uid"`
	P1     string `bson:"p1"`
	P2     string `bson:"p2"`
	Winner string `bson:"winner"`
}

type UserStat struct {
	TotalGamesPlayed int
	GamesWon         int
	GamesLost        int
	WinRate          float64
}

func processPlayer(uid string, winnerUID string, userStats map[string]*UserStat) {
	if uid == "" {
		return
	}

	stats, exists := userStats[uid]
	if !exists {
		stats = &UserStat{}
		userStats[uid] = stats
	}

	stats.TotalGamesPlayed++

	if uid == winnerUID {
		stats.GamesWon++
	} else {
		stats.GamesLost++
	}

	stats.WinRate = (float64(stats.GamesWon) / float64(stats.TotalGamesPlayed)) * 100
}
