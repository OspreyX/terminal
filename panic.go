package main

import (
	"io"
	"os"
	"strings"

	"github.com/mitchellh/panicwrap"
)

// This is output if a panic happens.
const panicOutput = `

!!!!!!!!!!!!!!!!!!!!!!!!!!! TERMINAL CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!

Terminal crashed! This is always indicative of a bug within Terminal.
A crash log has been placed at "crash.log" relative to your current
working directory. It would be immensely helpful if you could please
report the crash with Terminal[1] so that we can fix this.

[1]: https://github.com/intuition-io/terminal/issues

!!!!!!!!!!!!!!!!!!!!!!!!!!! TERMINAL CRASH !!!!!!!!!!!!!!!!!!!!!!!!!!!!
`

// panicHandler is what is called by panicwrap when a panic is encountered
// within Packer. It is guaranteed to run after the resulting process has
// exited so we can take the log file, add in the panic, and store it
// somewhere locally.
func panicHandler(logF *os.File) panicwrap.HandlerFunc {
	return func(m string) {
		// Write away just output this thing on stderr so that it gets
		// shown in case anything below fails.
		Log.Error("%s", m)

		// Create the crash log file where we'll write the logs
		f, err := os.Create("crash.log")
		if err != nil {
			Log.Error("Failed to create crash log file: %s", err)
			return
		}
		defer f.Close()

		// Seek the log file back to the beginning
		if _, err = logF.Seek(0, 0); err != nil {
			Log.Error("Failed to seek log file for crash: %s", err)
			return
		}

		// Copy the contents to the crash file. This will include
		// the panic that just happened.
		if _, err = io.Copy(f, logF); err != nil {
			Log.Error("Failed to write crash log: %s", err)
			return
		}

		// Tell the user a crash occurred in some helpful way that
		// they'll hopefully notice.
		Log.Info("\n")
		Log.Info(strings.TrimSpace(panicOutput))
	}
}
