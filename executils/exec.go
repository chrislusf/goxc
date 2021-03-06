// Utilities for invoking exec. Mainly focused on 'go build' and cross-compilation
package executils

/*
   Copyright 2013 Am Laher

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"io"
	"log"
	"os"
	"os/exec"
	//Tip for Forkers: please 'clone' from my url and then 'pull' from your url. That way you wont need to change the import path.
	//see https://groups.google.com/forum/?fromgroups=#!starred/golang-nuts/CY7o2aVNGZY
	"github.com/laher/goxc/platforms"
	"runtime"
	"strings"
)

// get list of args to be used in variable interpolation
// ldflags are used in any to any build-related go task (install,build,test)
func GetLdFlagVersionArgs(fullVersionName string) []string {
	if fullVersionName != "" {
		return []string{"-ldflags", "-X main.VERSION " + fullVersionName + ""}
	}
	return []string{}
}

// invoke the go command via the os/exec package
// 0.3.1
func InvokeGo(workingDirectory string, args []string, env []string, isVerbose bool) error {
	cmd := exec.Command("go")
	cmd.Args = append(cmd.Args, args...)
	cmd.Dir = workingDirectory
	f, err := RedirectIO(cmd)
	if err != nil {
		log.Printf("Error redirecting IO: %s", err)
	}
	if f != nil {
		defer f.Close()
	}
	//0.7.4 env replaces os.Environ
	cmd.Env = append(cmd.Env, env...)
	for _, thisProcessEnvItem := range os.Environ() {
		thisProcessEnvItemSplit := strings.Split(thisProcessEnvItem, "=")
		key := thisProcessEnvItemSplit[0]
		exists := false
		for _, specifiedEnvItem := range env {
			specifiedEnvItemSplit := strings.Split(specifiedEnvItem, "=")
			specifiedEnvKey := specifiedEnvItemSplit[0]
			if specifiedEnvKey == key {
				log.Printf("Overriding ENV variable (%s replaces %s)", specifiedEnvItem, thisProcessEnvItem)
				exists = true
			}
		}
		if !exists {
			cmd.Env = append(cmd.Env, thisProcessEnvItem)
		}
	}
	if isVerbose {
		log.Printf("(verbose!) all env vars for 'go': %s", cmd.Env)
	}
	if env != nil && len(env) > 0 {
		log.Printf("specified env vars for 'go': %s", env)
	}
	log.Printf("invoking 'go %v' from '%s'", PrintableArgs(args), workingDirectory)
	err = cmd.Start()
	if err != nil {
		log.Printf("Launch error: %s", err)
		return err
	} else {
		err = cmd.Wait()
		if err != nil {
			log.Printf("'go' returned error: %s", err)
			return err
		} else {
			if isVerbose {
				log.Printf("'go' completed successfully")
			}
		}
	}
	return nil
}

// returns a list of printable args
func PrintableArgs(args []string) string {
	ret := ""
	for _, arg := range args {
		if len(ret) > 0 {
			ret = ret + " "
		}
		if strings.Contains(arg, " ") {
			ret = ret + "\"" + arg + "\""
		} else {
			ret = ret + arg
		}
	}
	return ret
}

// this function copied from 'https://github.com/laher/mkdo'
// perhaps overkill for this project, but never mind.
func RedirectIO(cmd *exec.Cmd) (*os.File, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	//direct. Masked passwords work OK!
	cmd.Stdin = os.Stdin
	return nil, err
}

// check if cgoEnabled is required.
//0.2.4 refactored this out
// TODO not needed for go1.1+. Remove this once go1.0 reaches end of life. (when is that?)
func CgoEnabled(goos, arch string) string {
	var cgoEnabled string
	if goos == runtime.GOOS && arch == runtime.GOARCH {
		//note: added conditional in line with Dave Cheney, but this combination is not yet supported.
		if runtime.GOOS == platforms.FREEBSD && runtime.GOARCH == platforms.ARM {
			cgoEnabled = "0"
		} else {
			cgoEnabled = "1"
		}
	} else {
		cgoEnabled = "0"
	}
	return cgoEnabled
}
