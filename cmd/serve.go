// Copyright © 2018 Aurelien PERRIER <a.perrier89@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/perriea/tfstate/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a tfstate server",
	Long: `This command starts a tfstate server. By default, tfstate will start 
	and responds only API requests. You can start Ui with --ui flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run(true)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
