package cli

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Describe("splitCommand", func() {
		It("should split simple command", func() {
			cmd, args, err := splitCommand("vim file.txt")
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd).To(Equal("vim"))
			Expect(args).To(Equal([]string{"file.txt"}))
		})

		It("should handle command with multiple arguments", func() {
			cmd, args, err := splitCommand("git commit -m \"test message\"")
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd).To(Equal("git"))
			Expect(args).To(Equal([]string{"commit", "-m", "test message"}))
		})

		It("should handle command with no arguments", func() {
			cmd, args, err := splitCommand("ls")
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd).To(Equal("ls"))
			Expect(args).To(BeEmpty())
		})

		It("should handle quoted arguments", func() {
			cmd, args, err := splitCommand("echo \"hello world\" 'single quotes'")
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd).To(Equal("echo"))
			Expect(args).To(Equal([]string{"hello world", "single quotes"}))
		})

		It("should handle escaped quotes", func() {
			cmd, args, err := splitCommand(`echo "hello \"world\""`)
			Expect(err).ToNot(HaveOccurred())
			Expect(cmd).To(Equal("echo"))
			Expect(args).To(Equal([]string{`hello "world"`}))
		})

		It("should return error for invalid quotes", func() {
			_, _, err := splitCommand(`echo "unclosed quote`)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for empty string", func() {
			_, _, err := splitCommand("")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("editFile", func() {
		var originalEditor string

		BeforeEach(func() {
			originalEditor = os.Getenv("EDITOR")
		})

		AfterEach(func() {
			os.Setenv("EDITOR", originalEditor)
		})

		It("should return error when EDITOR is not set", func() {
			os.Unsetenv("EDITOR")
			err := editFile("test.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("set EDITOR environment variable"))
		})

		It("should return error for invalid EDITOR format", func() {
			os.Setenv("EDITOR", `vim "unclosed quote`)
			err := editFile("test.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("could not parse EDITOR"))
		})

		It("should return error when editor command fails", func() {
			// Use a command that doesn't exist
			os.Setenv("EDITOR", "nonexistenteditor")
			err := editFile("test.txt")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("could not edit file"))
		})

		// Note: Testing successful execution would require mocking exec.Command
		// which is beyond the scope of simple unit tests
	})
})