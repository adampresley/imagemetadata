# Image Metadata

This is a small library that pulls metadata from images. Metadata includes things like IPTC, EXIF and XMP data. It does this using other, well established libraries, then normalizes the data into a single structure.

## 🚀 Getting Started

First add the library as a dependency.

```bash
go get github.com/adampresley/imagemetadata
```

Now, use it like so.

```go
package main

import (
  "fmt"
  "os"

  "github.com/adampresley/imagemetadata"
)

func main() {
  var (
    err       error
    f         *os.File
    imageData *imagemetadata.ImageData
  )
 
  if f, err = os.Open("test-image.jpg"); err != nil {
    panic(err)
  }
 
  defer f.Close()
 
  if imageData, err = imagemetadata.NewFromJPEG(f); err != nil {
    panic(err)
  }
 
  fmt.Printf("%s\n", imageData.String())

  // Output:
  // ❯ go run main.go
  // Title: Title
  // Make: SONY
  // Model: ILCE-7M4
  // Lens Make:
  // Lens Model: E 28-75mm F2.8 A063
  // Caption: EXIF Description
  // Created Date: 1961-11-28T11:42:49
  // Width: 7000
  // Height: 4667
  // Latitude: 28.5378388
  // Longitude: -75.178244
  // Keywords: 2024, People, Thing
  // People: Person 1, Person 2
}
```
