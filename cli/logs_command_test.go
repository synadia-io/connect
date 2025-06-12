package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogsCommand", func() {
	var (
		cmd    *logsCommand
		appCtx *AppContext
	)

	BeforeEach(func() {
		appCtx, _ = newMockAppContext()
		
		cmd = &logsCommand{
			opts: &Options{},
		}
	})

	Describe("logs", func() {
		It("should capture logs for all connectors", func() {
			// The logs command doesn't have specific connector or instance parameters
			// It streams logs from all connectors until interrupted
			err := cmd.logsWithClient(appCtx)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})