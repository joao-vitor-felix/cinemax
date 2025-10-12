package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joao-vitor-felix/cinemax/internal/database"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	var port int
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()
	db := database.OpenPool()
	//TODO: move i18n setup to a separate function/file
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("internal/locales/en.json")
	bundle.LoadMessageFile("internal/locales/pt-BR.json")
	localizer := i18n.NewLocalizer(bundle, language.English.String(), language.BrazilianPortuguese.String())
	fmt.Println(localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "hello",
	}))

	//TODO: setup routes

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Handler: router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	log.Fatal(srv.ListenAndServe())
}
