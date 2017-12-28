package integration_test

import (
	"io"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var goodOutput = `data,duration(ms)
1d0afffafffafffafffafffafffa0cfffafffafffafffafffafffa01fffa0000,3250
1efffafffafffafffafffa03fffafffafffafffafffafffa05fffafffafffa00,3500
1efffafffafffa06fffafffafffafffafffafffa08fffafffafffafffafffa00,3500
02fffa0000000000000000000000000000000000000000000000000000000000,250
`
var _ = Describe("Performgen Integration", func() {
	var (
		stdin  io.WriteCloser
		stdout *gbytes.Buffer
		stderr *gbytes.Buffer
		cmd    *exec.Cmd
	)
	BeforeEach(func() {
		cmd = exec.Command(binaryPath)
		var err error
		stdin, err = cmd.StdinPipe()
		Expect(err).ToNot(HaveOccurred())
		stdout = gbytes.NewBuffer()
		stderr = gbytes.NewBuffer()
		cmd.Stdout = stdout
		cmd.Stderr = stderr

		err = cmd.Start()
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		cmd.Process.Kill()
	})
	It("converts MML to perform blocks as comma separated values", func(done Done) {
		_, err := stdin.Write([]byte("t80o3a2b2c2d2e2f2g2"))
		Expect(err).ToNot(HaveOccurred())
		Expect(stdin.Close()).To(Succeed())

		Expect(cmd.Wait()).To(Succeed())
		Expect(string(stdout.Contents())).To(Equal(goodOutput))
		Expect(string(stderr.Contents())).To(BeEmpty())
		close(done)
	}, 1.5)
	It("errors to stderr if the input is invalid", func(done Done) {
		_, err := stdin.Write([]byte("t80o6a2b2c2d2e2f2g2"))
		Expect(err).ToNot(HaveOccurred())
		Expect(stdin.Close()).To(Succeed())

		err = cmd.Wait()
		Expect(err).To(MatchError(ContainSubstring("exit status 1")))
		Expect(string(stdout.Contents())).To(BeEmpty())
		Expect(string(stderr.Contents())).To(ContainSubstring("execution error"))
		close(done)
	}, 1.5)
})
