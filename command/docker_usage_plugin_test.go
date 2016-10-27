/**
 * Copyright 2016 ECS Team, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package command_test

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"

	. "github.com/ecsteam/docker-usage/command"

	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	pluginio "github.com/ecsteam/cfcli-plugin-utils/plugin/io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Docker-Usage", func() {
	var fakeCliConnection *pluginfakes.FakeCliConnection

	var convertCommandOutputToStringSlice func(cmd *DockerUsagePlugin) []string

	var decolorizerRegex = regexp.MustCompile(`\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[m|K]`)

	var decolorize func(message string) string

	BeforeEach(func() {
		fakeCliConnection = &pluginfakes.FakeCliConnection{}
		fakeCliConnection.CliCommandWithoutTerminalOutputStub = func(args ...string) ([]string, error) {
			var output []string

			fixtureName := "fixtures" + args[1] + ".json"

			file, err := os.Open(fixtureName)
			defer file.Close()
			if err != nil {
				Fail("Could not open " + fixtureName)
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				output = append(output, scanner.Text())
			}

			return output, scanner.Err()
		}

		decolorize = func(message string) string {
			return string(decolorizerRegex.ReplaceAll([]byte(message), []byte("")))
		}

		convertCommandOutputToStringSlice = func(cmd *DockerUsagePlugin) []string {
			var lines []string
			scanner := bufio.NewScanner(bytes.NewBuffer(cmd.UI.Writer().(*bytes.Buffer).Bytes()))
			for scanner.Scan() {
				line := scanner.Text()
				lines = append(lines, line)
			}

			return lines
		}
	})

	Describe("Found Some Docker Images", func() {
		var output io.Writer
		var input io.Reader

		var cmd *DockerUsagePlugin
		BeforeEach(func() {
			output = new(bytes.Buffer)
			input = new(bytes.Buffer)

			ui := pluginio.UI{
				Input:  input,
				Output: output,
			}

			cmd = NewPlugin(ui)
		})

		It("returns some apps", func() {
			cmd.Run(fakeCliConnection, []string{"docker-usage"})

			lines := convertCommandOutputToStringSlice(cmd)

			Ω(len(lines)).To(Equal(7))
			Ω(decolorize(lines[1])).Should(BeEmpty())
			Ω(decolorize(lines[2])).Should(Equal("OK"))
			Ω(decolorize(lines[3])).Should(BeEmpty())
			Ω(decolorize(lines[4])).Should(MatchRegexp("image\\s*org\\s*space\\s*application"))
			Ω(decolorize(lines[5])).Should(HavePrefix("docker-repository.mydomain.com:5000/test/firehose-to-loginsight:v1"))
			Ω(decolorize(lines[6])).Should(HavePrefix("jghiloni/boot-test:latest"))
		})
	})
})
