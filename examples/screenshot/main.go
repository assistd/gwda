package main

import (
	"github.com/electricbubble/gwda"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	capture := func(name string) (int, error) {
		screenshot, err := driver.Screenshot()
		if err != nil {
			log.Fatal(err)
		}

		buf, _ := ioutil.ReadAll(screenshot)
		err = os.WriteFile(name, buf, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return len(buf), err
	}

	settings, err := driver.GetAppiumSettings()
	log.Printf("get settings: %#v", settings)

	// set wda screenshot quality to high
	_, err = driver.SetAppiumSettings(map[string]interface{}{
		"screenshotQuality":1,
	})
	t1 := time.Now()
	length, err := capture("full_quality.png")
	elapsed := time.Since(t1)
	log.Printf("full quality time:%v, size:%v", elapsed, length)

	// set wda screenshot quality to low
	_, err = driver.SetAppiumSettings(map[string]interface{}{
		"screenshotQuality":2,
	})
	t1 = time.Now()
	length, err = capture("low_quality.png")
	elapsed = time.Since(t1)
	log.Printf("low quality time:%v, size:%v", elapsed, length)

	/*
	img, format, err := image.Decode(screenshot)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = img, format
	userHomeDir, _ := os.UserHomeDir()
	file, err := os.Create(userHomeDir + "/Desktop/s1." + format)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = file.Close() }()
	switch format {
	case "png":
		err = png.Encode(file, img)
	case "jpeg":
		err = jpeg.Encode(file, img, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println(file.Name())
	 */
}
