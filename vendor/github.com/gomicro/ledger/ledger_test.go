package ledger

import (
	"fmt"
	"testing"

	"github.com/gomicro/penname"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestLogging(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Logging", func() {
		g.It("should write log lines", func() {
			mw := penname.New()
			l := New(mw, DebugLevel)
			msg := "some sorta message"

			for lvl := FatalLevel; lvl <= DebugLevel; lvl++ {
				l.write(lvl, msg)
				Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
			}
		})

		g.It("should honor the threshold", func() {
			mw := penname.New()
			msg := "some sorta message"

			l := New(mw, DebugLevel)
			for lvl := FatalLevel; lvl <= DebugLevel; lvl++ {
				l.write(lvl, msg)
				Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
			}

			l = New(mw, InfoLevel)
			for lvl := FatalLevel; lvl <= InfoLevel; lvl++ {
				l.write(lvl, msg)
				Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				mw.Reset()
			}
			l.write(DebugLevel, msg)
			Expect(string(mw.Written)).NotTo(Equal(fmt.Sprintf("%s: [%v]\n", DebugLevel, msg)))
			Expect(string(mw.Written)).To(Equal(""))

			l = New(mw, WarnLevel)
			for lvl := FatalLevel; lvl <= WarnLevel; lvl++ {
				l.write(lvl, msg)
				Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				mw.Reset()
			}
			for lvl := DebugLevel; lvl > WarnLevel; lvl-- {
				l.write(lvl, msg)
				Expect(string(mw.Written)).NotTo(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				Expect(string(mw.Written)).To(Equal(""))
			}

			l = New(mw, ErrorLevel)
			for lvl := FatalLevel; lvl <= ErrorLevel; lvl++ {
				l.write(lvl, msg)
				Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				mw.Reset()
			}
			for lvl := DebugLevel; lvl > ErrorLevel; lvl-- {
				l.write(lvl, msg)
				Expect(string(mw.Written)).NotTo(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				Expect(string(mw.Written)).To(Equal(""))
			}

			l = New(mw, FatalLevel)
			l.write(FatalLevel, msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", FatalLevel, msg)))
			mw.Reset()
			for lvl := DebugLevel; lvl > FatalLevel; lvl-- {
				l.write(lvl, msg)
				Expect(string(mw.Written)).NotTo(Equal(fmt.Sprintf("%s: [%v]\n", lvl, msg)))
				Expect(string(mw.Written)).To(Equal(""))
			}
		})

		g.It("should log levels", func() {
			mw := penname.New()
			l := New(mw, DebugLevel)
			msg := "bad wolf"

			l.Debug(msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", DebugLevel, msg)))

			l.Info(msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", InfoLevel, msg)))

			l.Warn(msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", WarnLevel, msg)))

			l.Error(msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", ErrorLevel, msg)))

			l.Fatal(msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v]\n", FatalLevel, msg)))
		})

		g.It("should log formatted levels", func() {
			mw := penname.New()
			l := New(mw, DebugLevel)
			msg := "bad wolf"

			l.Debugf("Additional Info: %v", msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [Additional Info: %v]\n", DebugLevel, msg)))

			l.Infof("Additional Info: %v", msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [Additional Info: %v]\n", InfoLevel, msg)))

			l.Warnf("Additional Info: %v", msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [Additional Info: %v]\n", WarnLevel, msg)))

			l.Errorf("Additional Info: %v", msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [Additional Info: %v]\n", ErrorLevel, msg)))

			l.Fatalf("Additional Info: %v", msg)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [Additional Info: %v]\n", FatalLevel, msg)))
		})

		g.It("should log multiple items", func() {
			mw := penname.New()
			l := New(mw, DebugLevel)
			msg1 := "exterminate"
			msg2 := "doctor!"

			l.Debug(msg1, msg2)
			Expect(string(mw.Written)).To(Equal(fmt.Sprintf("%s: [%v] [%v]\n", DebugLevel, msg1, msg2)))
		})
	})
}
