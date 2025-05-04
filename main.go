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

	"rss-score/api"
	"rss-score/db"
	"rss-score/service"

	_ "modernc.org/sqlite"
)

var (
	vaultDBEnvName  = "ZETTELKASTEN_DB"
	vaultCmd        = "pass"
	apiKeyVaultPath = "dev/rss-score/api-key"

	scoreFlagName = "score"
	idFlagName    = "id"
)

func main() {
	// validate options
	score := flag.Int(scoreFlagName, 0, "video score")
	videoID := flag.String(idFlagName, "", "expects YouTube video id")

	flag.Parse()

	for _, name := range []string{scoreFlagName, idFlagName} {
		checkTrue(isFlagPassed(name), "%s flag is unset", name)
	}

	checkTrue(len(*videoID) == 11, "expects len of video id to be 11, but got %d", len(*videoID))

	// validate environment
	dbPathZettelkastenPath, ok := os.LookupEnv(vaultDBEnvName)
	checkTrue(ok, "Missing %s", vaultDBEnvName)

	_, err := exec.LookPath(vaultCmd)
	checkNoErr(err)

	// setup service
	sqlite, err := sql.Open("sqlite", dbPathZettelkastenPath)
	checkNoErr(err)
	defer sqlite.Close()
	store := db.New(sqlite)

	apiKey, err := exec.Command(vaultCmd, apiKeyVaultPath).Output()
	checkTrue(!slices.Equal(apiKey, []byte{}), "api-key is empty!")
	checkNoErr(err)
	api := api.New(apiKey)

	svc := service.New(api, store)

	// process
	checkNoErr(svc.Run(*videoID, *score))
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
