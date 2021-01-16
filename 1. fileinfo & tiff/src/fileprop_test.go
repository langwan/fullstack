package passtime

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestProps(t *testing.T) {
	props, moreProps, _ := Props("/Users/langwan/Documents/data/passtime/source/tests/simples/file.cr2")
	for _, prop := range props {
		log.Info().Interface("prop", prop).Send()
	}

	for _, prop := range moreProps {
		log.Info().Interface("prop", prop).Send()
	}
}
