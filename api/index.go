package api

import (
    "bytes"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/juliencherry/local-datetime/geo"
    "github.com/juliencherry/local-datetime/textimage"
)

func Handler(w http.ResponseWriter, r *http.Request) {

   queryValues, ok := r.URL.Query()["t"]
   if !ok || len(queryValues) != 1 {
      log.Println("Unable to parse query")
      return
   }

   datetimeStr := strings.Replace(queryValues[0], "%20", " ", -1)
   datetime, err := time.Parse("2006-01-02 15:04 -0700", datetimeStr)
   if (err != nil) {
      log.Println("Unable to parse query")
      return
   }

   ip := r.Header.Get("X-FORWARDED-FOR")
   if ip == "" {
      ip = r.RemoteAddr
   }

   geolocator := geo.NewLocator(ip)

   location, err := time.LoadLocation(geolocator.Timezone())
   if (err == nil) {
      datetime = datetime.In(location)
   }

   imageText := []string{
      datetime.Format("Monday, January 2"),
      datetime.Format("3:04 PM") + " " + geolocator.City(),
   }

   imageBuffer := new(bytes.Buffer)

   err = textimage.Write(imageText, imageBuffer)
   if (err != nil) {
		log.Println(err)
      return
   }

   w.Header().Set("Content-Type", "image/png")
   w.Header().Set("Content-Length", strconv.Itoa(len(imageBuffer.Bytes())))
   if _, err := w.Write(imageBuffer.Bytes()); err != nil {
      log.Println("Unable to write image")
   }
}
