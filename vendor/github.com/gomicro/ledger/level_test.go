package ledger

import (
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestLevels(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Levels", func() {
		g.Describe("Stringing", func() {
			g.It("should output a level as a string", func() {
				testData := []struct {
					in  Level
					out string
				}{
					{FatalLevel, "FATAL"},
					{ErrorLevel, "ERROR"},
					{WarnLevel, "WARN"},
					{InfoLevel, "INFO"},
					{DebugLevel, "DEBUG"},
				}

				for _, r := range testData {
					i := r.in.String()
					Expect(i).To(Equal(r.out))
				}
			})
		})

		g.Describe("Parsing", func() {
			g.It("should parse a string to a level", func() {
				testData := []struct {
					in  string
					out Level
				}{
					{"FATAL", FatalLevel},
					{"ERROR", ErrorLevel},
					{"WARN", WarnLevel},
					{"INFO", InfoLevel},
					{"DEBUG", DebugLevel},
				}

				for _, r := range testData {
					i := ParseLevel(r.in)
					Expect(i).To(Equal(r.out))
				}
			})
		})
	})
}
