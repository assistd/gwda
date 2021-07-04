package main

import (
	"github.com/electricbubble/gwda"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
)

func loopJpegServer(client *http.Client) error {
	req, err := http.NewRequest(http.MethodGet, "http://localhost", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	mediaType, params, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	log.Printf("mediatype:%v, params:%#v, err:%v", mediaType, params, err)
	if err != nil {
		return err
	}

	mr := multipart.NewReader(resp.Body, params["boundary"][2:])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("error: reading next part", err)
			return err
		}

		jp, err := ioutil.ReadAll(part)
		if err != nil {
			log.Println("error: reading from bytes ", err)
			continue
		}
		log.Printf("len: %d", len(jp))
	}

	//for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
	//	value, _ := ioutil.ReadAll(part)
	//	log.Printf("len: %d", len(value))
	//}
	return nil
}

func main() {
	driver, err := gwda.NewUSBDriver(nil)
	if err != nil {
		log.Fatalln(err)
	}

	settings, err := driver.SetAppiumSettings(map[string]interface{}{
		"mjpegServerFramerate": 1,  //int, value:1~60, fps
		"mjpegScalingFactor":   50, //int, value:1~100, scale factor
		"screenshotQuality":    2, //int, value:1 or 2, 2 means low quality, see WDA object-c source code
	})
	log.Printf("settings: %#v", settings)

	err = loopJpegServer(driver.GetMjpegHTTPClient())
	if err != nil {
		log.Fatalf("loop error: %v", err)
	}
}
