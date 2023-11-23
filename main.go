package main

import (
	"fmt"
	"github.com/apognu/gocal"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	toAddresses []string
	authString  string
	hostString  string
)

func main() {
	// Configure logrus
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Hello World!")

	// Read environment variables
	toAddresses = strings.Split(os.Getenv("TO_EMAIL"), ",")
	authString = os.Getenv("NTFY_AUTH")
	hostString = os.Getenv("NTFY_HOST")

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Current Working Directory is = %s", dir)

	// Read and split multiple folders
	folderList := os.Getenv("ICS_DIR")
	if folderList == "" || len(toAddresses) == 0 || authString == "" || hostString == "" {
		log.Fatal("Es fehlen noch einige Parameter!!!\nICS_DIR, TO_EMAIL, NTFY_AUTH, NTFY_HOST")
	}

	folders := strings.Split(folderList, ",")
	for _, folder := range folders {
		log.Infof("folder: %s to: %s password: %s host: %s", folder, toAddresses, authString, hostString)
		listFilesForFolder(folder)
	}
}

func listFilesForFolder(folder string) {
	// Öffne das Verzeichnis
	dir, err := os.Open(folder)
	if err != nil {
		log.Fatal("Fehler beim Öffnen des Ordners:", err)
		return
	}
	defer dir.Close()

	// Lies alle Dateien im Verzeichnis
	dateien, err := dir.Readdir(0)
	if err != nil {
		log.Fatal("Fehler beim Lesen des Verzeichnisses:", err)
		return
	}

	// Durchlaufe die Liste der Dateien und gebe ihre Namen aus
	for _, datei := range dateien {
		// Überprüfe, ob es sich um ein Verzeichnis handelt. Wenn ja, ignoriere es.
		if datei.IsDir() {
			continue
		}

		// Hier kannst du die Dateinamen ausgeben oder damit arbeiten.
		log.Println(datei.Name())

		getNotifications(folder + "/" + datei.Name())
	}
}

func getNotifications(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tz, _ := time.LoadLocation("Europe/Berlin")

	start, end := truncateToDay(time.Now(), tz), truncateToDay(time.Now(), tz).Add(24*60*time.Minute)

	c := gocal.NewParser(f)
	c.Start, c.End = &start, &end
	c.Parse()

	for _, e := range c.Events {
		log.Infof("%s on %s", e.Summary, e.Start)

		messageSubject := fmt.Sprintf("Es existiert für heute ein neuer Kalendereintrag Namens: %s", e.Summary)
		log.Println(messageSubject)
		messageText := fmt.Sprintf("Der Termin beginnt heute um: %s und endet um: %s.", e.Start, e.End)
		if len(e.Location) != 0 {
			messageText += fmt.Sprintf("\n\nEr findet in %s statt.", e.Location)
		}
		if len(e.Description) != 0 {
			messageText += fmt.Sprintf("\n\nFolgende Notiz existiert in diesen Eintrag: \n%s", e.Description)
		}
		messageText += "\n\n This message is a service from go-ical-ntfy-reminder Version 1.1 written in Golang. \n Delivered by Simon Rieger"

		for _, toAddress := range toAddresses {
			sendMessage(messageSubject, messageText, toAddress)
		}
	}
}

func truncateToDay(t time.Time, tz *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz)
}

func sendMessage(messageSubject string, messageText string, toAddress string) {

	req, _ := http.NewRequest("POST", hostString,
		strings.NewReader(messageText))
	req.Header.Set("Title", messageSubject)
	req.Header.Set("Authorization", "Basic "+authString)
	req.Header.Set("Email", toAddress)
	req.Header.Set("Tags", "date,octopus")
	req.Header.Set("Priority", "high")
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Ntfy Message not Sent Successfully, Error: " + err.Error())
		return
	}

	log.Infof("Ntfy Message Sent Successfully, Status: " + do.Status)
}
