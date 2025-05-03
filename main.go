// dependencies:
// - sqlite база данных хранилища со схемой от 2025-05-03
// - команда pass и ключ к апи, сохраннённый в нём
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/exec"
	"slices"

	_ "modernc.org/sqlite"
	"rss-score/api"
	"rss-score/db"
)

var (
	vaultDBEnvName  = "ZETTELKASTEN_DB"
	apiKeyVaultPath = "dev/rss-score/api-key"

	scoreFlagName = "score"
	idFlagName    = "id"
)

func main() {
	// validate environment
	dbPathZettelkasten, ok := os.LookupEnv(vaultDBEnvName)
	checkTrue(ok, "Missing %s", vaultDBEnvName)

	conn, err := sql.Open("sqlite", dbPathZettelkasten)
	checkNoErr(err)
	defer conn.Close()

	apiKey, err := exec.Command("pass", apiKeyVaultPath).Output()
	checkNoErr(err)
	checkTrue(!slices.Equal(apiKey, []byte{}), "api-key is empty!")
	api := api.New(apiKey)

	// validate options
	score := flag.Int(scoreFlagName, 0, "video score")
	videoID := flag.String(idFlagName, "", "expects YouTube video id")

	flag.Parse()

	for _, name := range []string{scoreFlagName, idFlagName} {
		checkTrue(isFlagPassed(name), "%s flag is unset", name)
	}

	// process
	article, err := api.FetchMetadata(*videoID)
	checkNoErr(err)

	checkNoErr(db.InsertOrUpdate(conn, article, *score))
	log.Println("success!")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func checkNoErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func checkTrue(ok bool, message string, args ...any) {
	if !ok {
		log.Fatalf(message, args...)
	}
}
