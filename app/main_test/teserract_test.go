package unittest

import (
	"fmt"
	"testing"

	"github.com/otiai10/gosseract/v2"
)

/*
brew install tesseract leptonica
go get -u github.com/otiai10/gosseract

after install and import gosseract make export on path project folder

export LIBRARY_PATH="/opt/homebrew/lib"
export CPATH="/opt/homebrew/include"
*/
func TestImageToText(t *testing.T) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("images/plat_no.jpeg")
	text, _ := client.Text()
	fmt.Println("text : ", text)
}
