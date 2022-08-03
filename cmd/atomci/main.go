/*
Copyright 2021 The AtomCI Group Authors.

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
	"runtime"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql" // import your used driver

	"github.com/go-atomci/atomci/internal/initialize"
	"github.com/go-atomci/atomci/internal/migrations"
	"github.com/go-atomci/atomci/internal/models"
	"github.com/go-atomci/atomci/version"

	"github.com/go-atomci/atomci/internal/cronjob"
	"github.com/go-atomci/atomci/internal/routers"
)

func main() {
	models.InitDB()
	migrations.Migrate()
	// TODO: resource items migrate later
	initialize.Init()

	cronjob.RunPublishJobServer()

	routers.RegisterRoutes()
	beego.Info("Beego version:", beego.VERSION)
	beego.Info("Golang version:", runtime.Version())
	version.PrintFullVersionInfo()
	beego.Run()
}
