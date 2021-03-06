package depot

import (
	"os"

	"github.com/cloudfoundry-incubator/executor/api"
	"github.com/cloudfoundry-incubator/executor/registry"
	"github.com/cloudfoundry-incubator/executor/sequence"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
)

type RunSequence struct {
	CompleteURL  string
	Registration api.Container
	Sequence     sequence.Step
	Result       *string
	Registry     registry.Registry
	Logger       lager.Logger
}

func (r RunSequence) Run(sigChan <-chan os.Signal, readyChan chan<- struct{}) error {
	seqComplete := make(chan error)

	runLog := r.Logger.Session("run", lager.Data{
		"guid":   r.Registration.Guid,
		"handle": r.Registration.ContainerHandle,
	})

	go func() {
		runLog.Info("starting")
		seqComplete <- r.Sequence.Perform()
	}()

	close(readyChan)

	for {
		select {
		case <-sigChan:
			sigChan = nil
			r.Sequence.Cancel()
			runLog.Info("cancelled")

		case err := <-seqComplete:
			if err == sequence.CancelledError {
				return err
			}

			runLog.Info("completed")

			payload := api.ContainerRunResult{
				Guid:   r.Registration.Guid,
				Result: *r.Result,
			}

			if err != nil {
				payload.Failed = true
				payload.FailureReason = err.Error()
			}

			err = r.Registry.Complete(r.Registration.Guid, payload)
			if err != nil {
				runLog.Error("failed-to-complete", err)
			}

			if r.CompleteURL == "" {
				return err
			}

			ifrit.Envoke(&Callback{
				URL:     r.CompleteURL,
				Payload: payload,
			})

			runLog.Info("callback-started")

			return err
		}
	}
}
