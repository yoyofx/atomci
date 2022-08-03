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

package api

import (
	"github.com/go-atomci/atomci/internal/core/apps"
	"github.com/go-atomci/atomci/internal/core/kuberes"
	"github.com/go-atomci/atomci/internal/middleware/log"
	"github.com/go-atomci/atomci/utils/errors"
)

// AppController ...
type AppController struct {
	BaseController
}

// CreateSCMApp for project
func (a *AppController) CreateSCMApp() {
	req := &apps.ScmAppReq{}
	a.DecodeJSONReq(&req)
	mgr := apps.NewAppManager()
	_, result := mgr.CreateSCMApp(req, a.User)
	if result != nil {
		a.HandleInternalServerError(result.Error())
		log.Log.Error("add project app error: %s", result.Error())
		return
	}
	a.Data["json"] = NewResult(true, result, "")
	a.ServeJSON()
}

// VerifySCMAppConnetion
// 验证仓库地址是否能连通
func (a *AppController) VerifySCMAppConnetion() {
	req := &apps.ScmAppReq{}
	a.DecodeJSONReq(&req)
	app := apps.NewAppManager()
	err := app.VerifyAppConnetion(req.RepoID, req.Path, req.FullName)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("verify scm app connetion occur error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, "", "")
	a.ServeJSON()
}

func (a *AppController) GetAllApps() {
	mgr := apps.NewAppManager()
	// TODO: add app tag filter base on permisson
	result, err := mgr.GetScmApps()
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get scm apps error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, result, "")
	a.ServeJSON()
}

// GetAppsByPagination ..
func (a *AppController) GetAppsByPagination() {
	filterQuery := a.GetFilterQuery()
	mgr := apps.NewAppManager()
	result, err := mgr.GetScmAppsByPagination(filterQuery)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get scm app list error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, result, "")
	a.ServeJSON()
}

// ScmAppInfo ..
func (a *AppController) ScmAppInfo() {
	scmAppID, _ := a.GetInt64FromPath(":app_id")
	mgr := apps.NewAppManager()
	result, err := mgr.GetScmApp(scmAppID)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get scm app error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, result, "")
	a.ServeJSON()
}

// UpdateScmApp ..
func (a *AppController) UpdateScmApp() {
	scmAppID, _ := a.GetInt64FromPath(":app_id")
	req := &apps.ScmAppUpdateReq{}
	a.DecodeJSONReq(req)
	am := apps.NewAppManager()
	if err := am.UpdateProjectApp(scmAppID, req); err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("update scm app error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, nil, "")
	a.ServeJSON()
}

// DeleteScmApp ...
func (a *AppController) DeleteScmApp() {
	scmAppID, _ := a.GetInt64FromPath(":app_id")
	am := apps.NewAppManager()
	result := am.DeleteSCMApp(scmAppID)
	if result != nil {
		a.HandleInternalServerError(result.Error())
		log.Log.Error("delete scm app error: %s", result.Error())
		return
	}
	a.Data["json"] = NewResult(true, result, "")
	a.ServeJSON()
}

// GetArrange ...
func (a *AppController) GetArrange() {
	appID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	envID, err := a.GetInt64FromPath(":env_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid env id"))
		return
	}
	mgr := apps.NewAppManager()
	arrange, err := mgr.GetArrange(appID, envID)
	if err != nil {
		a.ServeError(err)
		return
	}
	a.ServeResult(NewResult(true, arrange, ""))
}

// SetArrange ...
func (a *AppController) SetArrange() {
	projectAppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid project app id"))
		return
	}
	arrangeEnvID, err := a.GetInt64FromPath(":env_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid project env id"))
		return
	}
	request := apps.AppArrangeReq{}
	a.DecodeJSONReq(&request)

	native := &kuberes.NativeTemplate{
		Template: request.Config,
	}

	if err := native.Validate(); err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("yaml parse error: %s", err.Error()))
		return
	}

	mgr := apps.NewAppManager()
	err = mgr.SetArrange(projectAppID, arrangeEnvID, &request)
	if err != nil {
		a.ServeError(err)
		return
	}
	a.ServeResult(NewResult(true, nil, ""))
}

func (a *AppController) ParseArrangeYaml() {
	request := apps.AppArrangConfig{}
	a.DecodeJSONReq(&request)

	native := &kuberes.NativeTemplate{
		Template: request.Config,
	}
	rsp, err := native.GetContainerImages()
	if err != nil {
		log.Log.Debug("get container images error: %s")
	}
	a.ServeResult(NewResult(true, rsp, ""))
}

// GetGitProjectsByRepoID ..
func (a *AppController) GetGitProjectsByRepoID() {
	repoID, _ := a.GetInt64FromPath(":repo_id")
	mgr := apps.NewAppManager()
	rsp, err := mgr.GetScmProjectsByRepoID(repoID)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("get repo's projects error: %s", err.Error())
		return
	}
	a.ServeResult(NewResult(true, rsp, ""))
}

// GetAppBranches ..
func (a *AppController) GetAppBranches() {
	AppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	filterQuery := a.GetFilterQuery()
	mgr := apps.NewAppManager()
	rsp, err := mgr.AppBranches(AppID, filterQuery)
	if err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("Get app list error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, rsp, "")
	a.ServeJSON()
}

// SyncAppBranches ..
func (a *AppController) SyncAppBranches() {
	AppID, err := a.GetInt64FromPath(":app_id")
	if err != nil {
		a.ServeError(errors.NewBadRequest().SetMessage("invalid app id"))
		return
	}
	mgr := apps.NewAppManager()
	if err := mgr.SyncAppBranches(AppID); err != nil {
		a.HandleInternalServerError(err.Error())
		log.Log.Error("sync app branches error: %s", err.Error())
		return
	}
	a.Data["json"] = NewResult(true, nil, "")
	a.ServeJSON()
}
