package serr

import (
	"log"
	"testing"

	"github.com/pkg/errors"
)

func Test(t *testing.T) {
	err := test1()
	// log.Print(Wrap(err))
	log.Printf("%+v", err)
}

func test1() error {
	err := test2()
	return Wrap(err)
}

func test2() error {
	err := test3()
	return Wrap(err)
}

func test3() error {
	return errors.New("test3")
}
