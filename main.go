package main

import (
	"crypto/tls"
	"fmt"
	"github.com/apognu/gocal"
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	fromString     string
	toAddresses    []string
	passwordString string
	hostString     string
	hostPortString int
)

func main() {
	// Configure logrus
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.Info("Hello World!")

	// Initialisiere err mit einem Fehlerwert
	var err error

	// Read environment variables
	fromString = os.Getenv("FROM_EMAIL")
	toAddresses = strings.Split(os.Getenv("TO_EMAIL"), ",")
	passwordString = os.Getenv("EMAIL_PASSWORD")
	hostString = os.Getenv("SMTP_HOST")
	hostPortString, err = strconv.Atoi(os.Getenv("SMTP_PORT"))

	// Überprüfung auf Fehler
	if err != nil {
		fmt.Println("Fehler bei der Umwandlung. Verwende Standardwert 587.")
		// Setze einen Standardwert, zum Beispiel 0
		hostPortString = 587
	}

	// Get the current working directory
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Current Working Directory is = %s", dir)

	// Read and split multiple folders
	folderList := os.Getenv("ICS_DIR")
	if folderList == "" || fromString == "" || len(toAddresses) == 0 || passwordString == "" || hostString == "" {
		log.Fatal("Es fehlen noch einige Parameter!!!\nICS_DIR, FROM_EMAIL, TO_EMAIL, EMAIL_PASSWORD, SMTP_HOST")
	}

	folders := strings.Split(folderList, ",")
	for _, folder := range folders {
		log.Infof("folder: %s from: %s to: %s password: %s host: %s", folder, fromString, toAddresses, passwordString, hostString)
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

	var tzMapping = map[string]string{
		"My Super Zone": "Europe/Berlin",
	}

	gocal.SetTZMapper(func(s string) (*time.Location, error) {
		if tzid, ok := tzMapping[s]; ok {
			return time.LoadLocation(tzid)
		}
		return nil, fmt.Errorf("")
	})

	start, end := truncateToDay(time.Now()), truncateToDay(time.Now()).Add(24*60*time.Minute)

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
		messageText += "\n\n This email is a service from mail-reminder Version 1.0 written in Golang. \n Delivered by Simon Rieger"

		for _, toAddress := range toAddresses {
			sendMail(messageSubject, messageText, toAddress)
		}
	}
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func sendMail(messageSubject string, messageText string, toAddress string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", fromString)

	// Set E-Mail receivers
	m.SetHeader("To", toAddress)

	// Set E-Mail subject
	m.SetHeader("Subject", messageSubject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", messageText)

	// Settings for SMTP server
	d := gomail.NewDialer(hostString, hostPortString, fromString, passwordString)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Infof("Email Message Sent Successfully")
}
