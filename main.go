package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tebeka/selenium"
)

const (
	url             = "https://docs.google.com/forms/d/e/1FAIpQLSeYRQtIuz8nPe3yhUShXF6OJbZSjoJBBQQgfR7vBRfgwctHdw/viewform"
	seleniumPath    = "utils/selenium-server-standalone-3.141.59.jar"
	geckoDriverPath = "utils/geckodriver"
	port            = 8080
	times           = 100
	timeout         = time.Second * 2
)

func main() {

	stopCh := make(chan os.Signal, 2)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGQUIT)

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath),
	}

	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatal(err)
	}

	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))

	if err != nil {
		log.Fatal(err)
	}

	defer wd.Quit()

	for i := 0; i < times; i++ {
		var name string
		var email string

		switch i % 2 {
		case 0:
			name = "###"
			email = "###"
		case 1:
			name = "###"
			email = "###"

		}

		if err := wd.Get(url); err != nil {
			log.Fatal(err)
		}
		textFields, err := wd.FindElements(selenium.ByCSSSelector, ".quantumWizTextinputPaperinputInput")
		if err != nil {
			log.Fatal(err)
		}

		err = textFields[0].SendKeys(name)
		if err != nil {
			log.Fatal(err)
		}

		err = textFields[1].SendKeys(email)
		if err != nil {
			log.Fatal(err)
		}

		radios, err := wd.FindElements(selenium.ByCSSSelector, ".freebirdFormviewerViewItemsRadioOptionContainer")
		if err != nil {
			log.Fatal(err)
		}

		if err := radios[2].Click(); err != nil {
			log.Fatal(err)
		}

		if err := radios[16].Click(); err != nil {
			log.Fatal(err)
		}

		button, err := wd.FindElement(selenium.ByCSSSelector, ".quantumWizButtonPaperbuttonContent")
		if err != nil {
			log.Fatal(err)
		}

		if err := button.Click(); err != nil {
			log.Fatal(err)
		}

		time.Sleep(timeout)
	}

	for {
		select {
		case <-stopCh:
			os.Exit(0)
		}
	}

}
