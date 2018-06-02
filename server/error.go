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
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// NotFoundHandler : Handler 404
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	message := &ErrorCode{
		Error:   http.StatusNotFound,
		Message: "Not found",
	}
	resp, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(resp)
}
