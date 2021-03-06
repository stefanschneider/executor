package log_streamer_test

import (
	"fmt"
	"strings"
	"time"

	. "github.com/cloudfoundry-incubator/executor/log_streamer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/loggregatorlib/logmessage"
)

var _ = Describe("LogStreamer", func() {
	var loggregatorEmitter *FakeLoggregatorEmitter
	var streamer LogStreamer
	guid := "the-guid"
	sourceName := "the-source-name"
	index := 11

	BeforeEach(func() {
		loggregatorEmitter = NewFakeLoggregatorEmmitter()
		streamer = New(guid, sourceName, &index, loggregatorEmitter)
	})

	Context("when told to emit", func() {
		Context("when given a message that corresponds to one line", func() {
			BeforeEach(func() {
				fmt.Fprintln(streamer.Stdout(), "this is a log")
				fmt.Fprintln(streamer.Stdout(), "this is another log")
			})

			It("should emit that message", func() {
				Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))

				emission := loggregatorEmitter.Emissions[0]
				Ω(emission.GetAppId()).Should(Equal(guid))
				Ω(string(emission.GetMessage())).Should(Equal("this is a log"))
				Ω(emission.GetMessageType()).Should(Equal(logmessage.LogMessage_OUT))
				Ω(emission.GetSourceId()).Should(Equal("11"))

				emission = loggregatorEmitter.Emissions[1]
				Ω(emission.GetAppId()).Should(Equal(guid))
				Ω(emission.GetSourceName()).Should(Equal(sourceName))
				Ω(emission.GetSourceId()).Should(Equal("11"))
				Ω(string(emission.GetMessage())).Should(Equal("this is another log"))
				Ω(emission.GetMessageType()).Should(Equal(logmessage.LogMessage_OUT))
				Ω(*emission.Timestamp).Should(BeNumerically("~", time.Now().UnixNano(), 10*time.Millisecond))
			})
		})

		Context("when given a message with all sorts of fun newline characters", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "A\nB\rC\n\rD\r\nE\n\n\nF\r\r\rG\n\r\r\n\n\n\r")
			})

			It("should do the right thing", func() {
				Ω(loggregatorEmitter.Emissions).Should(HaveLen(7))
				for i, expectedString := range []string{"A", "B", "C", "D", "E", "F", "G"} {
					Ω(string(loggregatorEmitter.Emissions[i].GetMessage())).Should(Equal(expectedString))
				}
			})
		})

		Context("when given a series of short messages", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "this is a log")
				fmt.Fprintf(streamer.Stdout(), " it is made of wood")
				fmt.Fprintf(streamer.Stdout(), " - and it is longer")
				fmt.Fprintf(streamer.Stdout(), "than it seems\n")
			})

			It("concatenates them, until a new-line is received, and then emits that", func() {
				Ω(loggregatorEmitter.Emissions).Should(HaveLen(1))

				emission := loggregatorEmitter.Emissions[0]
				Ω(string(emission.GetMessage())).Should(Equal("this is a log it is made of wood - and it is longerthan it seems"))
			})
		})

		Context("when given a message with multiple new lines", func() {
			BeforeEach(func() {
				fmt.Fprintf(streamer.Stdout(), "this is a log\nand this is another\nand this one isn't done yet...")
			})

			It("should break the message up into multiple loggings", func() {
				Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))

				emission := loggregatorEmitter.Emissions[0]
				Ω(string(emission.GetMessage())).Should(Equal("this is a log"))

				emission = loggregatorEmitter.Emissions[1]
				Ω(string(emission.GetMessage())).Should(Equal("and this is another"))
			})
		})

		Describe("message limits", func() {
			var message string
			Context("when the message is just at the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", MAX_MESSAGE_SIZE)
					Ω([]byte(message)).Should(HaveLen(MAX_MESSAGE_SIZE), "Ensure that the byte representation of our message is under the limit")

					fmt.Fprintf(streamer.Stdout(), message)
				})

				It("should break the message up and send multiple messages", func() {
					Ω(loggregatorEmitter.Emissions).Should(HaveLen(1))
					emission := loggregatorEmitter.Emissions[0]
					Ω(string(emission.GetMessage())).Should(Equal(message))
				})
			})

			Context("when the message exceeds the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", MAX_MESSAGE_SIZE)
					message += strings.Repeat("8", MAX_MESSAGE_SIZE)
					message += strings.Repeat("9", MAX_MESSAGE_SIZE)
					message += "hello\n"
					fmt.Fprintf(streamer.Stdout(), message)
				})

				It("should break the message up and send multiple messages", func() {
					Ω(loggregatorEmitter.Emissions).Should(HaveLen(4))
					Ω(string(loggregatorEmitter.Emissions[0].GetMessage())).Should(Equal(strings.Repeat("7", MAX_MESSAGE_SIZE)))
					Ω(string(loggregatorEmitter.Emissions[1].GetMessage())).Should(Equal(strings.Repeat("8", MAX_MESSAGE_SIZE)))
					Ω(string(loggregatorEmitter.Emissions[2].GetMessage())).Should(Equal(strings.Repeat("9", MAX_MESSAGE_SIZE)))
					Ω(string(loggregatorEmitter.Emissions[3].GetMessage())).Should(Equal("hello"))
				})
			})

			Context("when having to deal with byte boundaries", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", MAX_MESSAGE_SIZE-1)
					message += "\u0623\n"
					fmt.Fprintf(streamer.Stdout(), message)
				})

				It("should break the message up and send multiple messages", func() {
					Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))
					Ω(string(loggregatorEmitter.Emissions[0].GetMessage())).Should(Equal(strings.Repeat("7", MAX_MESSAGE_SIZE-1)))
					Ω(string(loggregatorEmitter.Emissions[1].GetMessage())).Should(Equal("\u0623"))
				})
			})

			Context("while concatenating, if the message exceeds the emittable length", func() {
				BeforeEach(func() {
					message = strings.Repeat("7", MAX_MESSAGE_SIZE-2)
					fmt.Fprintf(streamer.Stdout(), message)
					fmt.Fprintf(streamer.Stdout(), "778888\n")
				})

				It("should break the message up and send multiple messages", func() {
					Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))
					Ω(string(loggregatorEmitter.Emissions[0].GetMessage())).Should(Equal(strings.Repeat("7", MAX_MESSAGE_SIZE)))
					Ω(string(loggregatorEmitter.Emissions[1].GetMessage())).Should(Equal("8888"))
				})
			})
		})
	})

	Context("when told to emit stderr", func() {
		It("should handle short messages", func() {
			fmt.Fprintf(streamer.Stderr(), "this is a log\nand this is another\nand this one isn't done yet...")
			Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))

			emission := loggregatorEmitter.Emissions[0]
			Ω(string(emission.GetMessage())).Should(Equal("this is a log"))
			Ω(emission.GetSourceName()).Should(Equal(sourceName))
			Ω(emission.GetMessageType()).Should(Equal(logmessage.LogMessage_ERR))

			emission = loggregatorEmitter.Emissions[1]
			Ω(string(emission.GetMessage())).Should(Equal("and this is another"))
		})

		It("should handle long messages", func() {
			fmt.Fprintf(streamer.Stderr(), strings.Repeat("e", MAX_MESSAGE_SIZE+1)+"\n")
			Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))

			emission := loggregatorEmitter.Emissions[0]
			Ω(string(emission.GetMessage())).Should(Equal(strings.Repeat("e", MAX_MESSAGE_SIZE)))

			emission = loggregatorEmitter.Emissions[1]
			Ω(string(emission.GetMessage())).Should(Equal("e"))
		})
	})

	Context("when told to flush", func() {
		It("should send whatever log is left in its buffer", func() {
			fmt.Fprintf(streamer.Stdout(), "this is a stdout")
			fmt.Fprintf(streamer.Stderr(), "this is a stderr")

			Ω(loggregatorEmitter.Emissions).Should(HaveLen(0))

			streamer.Flush()

			Ω(loggregatorEmitter.Emissions).Should(HaveLen(2))
			Ω(loggregatorEmitter.Emissions[0].GetMessageType()).Should(Equal(logmessage.LogMessage_OUT))
			Ω(loggregatorEmitter.Emissions[1].GetMessageType()).Should(Equal(logmessage.LogMessage_ERR))
		})
	})

	Context("when there is no app guid", func() {
		It("does nothing when told to emit or flush", func() {
			streamer = New("", sourceName, &index, loggregatorEmitter)

			streamer.Stdout().Write([]byte("hi"))
			streamer.Stderr().Write([]byte("hi"))
			streamer.Flush()

			Ω(loggregatorEmitter.Emissions).Should(BeEmpty())
		})
	})

	Context("when there is no source index", func() {
		It("defaults to 0", func() {
			streamer = New(guid, sourceName, nil, loggregatorEmitter)

			streamer.Stdout().Write([]byte("hi"))
			streamer.Flush()

			Ω(loggregatorEmitter.Emissions[0].GetSourceId()).Should(Equal("0"))
		})
	})
})

type FakeLoggregatorEmitter struct {
	Emissions []*logmessage.LogMessage
}

func NewFakeLoggregatorEmmitter() *FakeLoggregatorEmitter {
	return &FakeLoggregatorEmitter{}
}

func (e *FakeLoggregatorEmitter) Emit(appid, message string) {
	panic("no no no no")
}

func (e *FakeLoggregatorEmitter) EmitError(appid, message string) {
	panic("no no no no")
}

func (e *FakeLoggregatorEmitter) EmitLogMessage(msg *logmessage.LogMessage) {
	e.Emissions = append(e.Emissions, msg)
}
