package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"zxq.co/ripple/agplwarning"
	"github.com/kawatapw/api/app"
	"github.com/kawatapw/api/beatmapget"
	"github.com/kawatapw/api/common"
	// Golint pls dont break balls
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/serenize/snaker"
	"gopkg.in/thehowl/go-osuapi.v1"
)

// Version is the git hash of the application. Do not edit. This is
// automatically set using -ldflags during build time.
var Version string

func init() {
	log.SetFlags(log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d|", syscall.Getpid()))
	common.Version = Version
}

var db *sqlx.DB

func main() {
	err := agplwarning.Warn("ripple", "Ripple API")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("Kawata API")
	if Version != "" {
		fmt.Print("; git commit hash: ", Version)
	}
	fmt.Println()

	conf, halt := common.Load()
	if halt {
		return
	}

	if !strings.Contains(conf.DSN, "parseTime=true") {
		c := "?"
		if strings.Contains(conf.DSN, "?") {
			c = "&"
		}
		conf.DSN += c + "parseTime=true&charset=utf8mb4,utf8&collation=utf8mb4_general_ci"
	}

	db, err = sqlx.Open(conf.DatabaseType, conf.DSN)
	if err != nil {
		log.Fatalln(err)
	}

	db.MapperFunc(func(s string) string {
		if x, ok := commonClusterfucks[s]; ok {
			return x
		}
		return snaker.CamelToSnake(s)
	})

	beatmapget.Client = osuapi.NewClient(conf.OsuAPIKey)
	beatmapget.DB = db

	engine := app.Start(conf, db)

	startuato(engine.Handler)
}

var commonClusterfucks = map[string]string{
	"RegisteredOn": "register_datetime",
	"UsernameAKA":  "username_aka",
	"BeatmapMD5":   "beatmap_md5",
	"Count300":     "300_count",
	"Count100":     "100_count",
	"Count50":      "50_count",
	"CountGeki":    "gekis_count",
	"CountKatu":    "katus_count",
	"CountMiss":    "misses_count",
	"PP":           "pp",
}
