/*
Copyright © 2021 Leopold Fitz caojunkaiv@gmail.com
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
package main

import (
	"log"
	"strings"
	"time"

	"github.com/leagos/sniper-bot/cmd"
	"github.com/leagos/sniper-bot/utils"
)

func main() {
	cpuid := utils.GetCpuId()
	var license string = "178BFBFF00100F53"
	timeLayout := "2006-01-02 15:04:05"
	expireDate, _ := time.Parse(timeLayout, "2023-03-05 12:00:00")
	if time.Now().Unix() > expireDate.Unix() {
		log.Fatalf("Expired")
	}
	if !strings.Contains(cpuid, license) {
		log.Fatalf("No license")
	}
	cmd.Execute()
}
