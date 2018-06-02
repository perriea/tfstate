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

package server

import (
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func getState(w http.ResponseWriter, r *http.Request) {
	var (
		message, username, password string
	)

	// search in config file (user & password)

	if BasicAuth(w, r, username, password, "Provide user name and password") {
		message = r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Hello " + message
	} else {
		message = "test"
	}

	w.Write([]byte(message))
}
