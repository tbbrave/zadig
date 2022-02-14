/*
Copyright 2021 The KodeRover Authors.

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

package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/koderover/zadig/pkg/util/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Router struct {
	metricCollector       *metrics.Metric
	responseTimeHistogram prometheus.HistogramVec //in milliseconds
}

func NewRouter() *Router {
	metricCollector := &metrics.Metric{Type: metrics.Histogram}
	metricCollector.Name = "ResponseTime"
	metricCollector.Namespace = "Core"
	metricCollector.Labels = []string{"API"}
	metrics.RegisterMetric(metricCollector)

	return &Router{
		metricCollector:       metricCollector,
		responseTimeHistogram: metricCollector.Collector.((prometheus.HistogramVec))}
}

func (r *Router) Inject(router *gin.RouterGroup) {
	roles := router.Group("roles")
	{
		roles.POST("", r.metricHandlerWrapper(CreateRole, ""))
		roles.POST("/bulk-delete", r.metricHandlerWrapper(DeleteRoles, ""))
		roles.PATCH("/:name", r.metricHandlerWrapper(UpdateRole, ""))
		roles.PUT("/:name", r.metricHandlerWrapper(UpdateOrCreateRole, ""))
		roles.GET("", r.metricHandlerWrapper(ListRoles, ""))
		roles.GET("/:name", r.metricHandlerWrapper(GetRole, ""))
		roles.DELETE("/:name", r.metricHandlerWrapper(DeleteRole, ""))
	}

	publicRoles := router.Group("public-roles")
	{
		publicRoles.POST("", r.metricHandlerWrapper(CreatePublicRole, ""))
		publicRoles.GET("", r.metricHandlerWrapper(ListPublicRoles, ""))
		publicRoles.GET("/:name", r.metricHandlerWrapper(GetPublicRole, ""))
		publicRoles.PATCH("/:name", r.metricHandlerWrapper(UpdatePublicRole, ""))
		publicRoles.PUT("/:name", r.metricHandlerWrapper(UpdateOrCreatePublicRole, ""))
		publicRoles.DELETE("/:name", r.metricHandlerWrapper(DeletePublicRole, ""))
	}

	systemRoles := router.Group("system-roles")
	{
		systemRoles.POST("", r.metricHandlerWrapper(CreateSystemRole, ""))
		systemRoles.PUT("/:name", r.metricHandlerWrapper(UpdateOrCreateSystemRole, ""))
		systemRoles.GET("", r.metricHandlerWrapper(ListSystemRoles, ""))
		systemRoles.DELETE("/:name", r.metricHandlerWrapper(DeleteSystemRole, ""))
	}

	roleBindings := router.Group("rolebindings")
	{
		roleBindings.POST("", r.metricHandlerWrapper(CreateRoleBinding, ""))
		roleBindings.PUT("/:name", r.metricHandlerWrapper(UpdateRoleBinding, ""))
		roleBindings.GET("", r.metricHandlerWrapper(ListRoleBindings, ""))
		roleBindings.DELETE("/:name", r.metricHandlerWrapper(DeleteRoleBinding, ""))
		roleBindings.POST("/bulk-delete", r.metricHandlerWrapper(DeleteRoleBindings, ""))
	}

	systemRoleBindings := router.Group("system-rolebindings")
	{
		systemRoleBindings.POST("", r.metricHandlerWrapper(CreateSystemRoleBinding, ""))
		systemRoleBindings.GET("", r.metricHandlerWrapper(ListSystemRoleBindings, ""))
		systemRoleBindings.DELETE("/:name", r.metricHandlerWrapper(DeleteSystemRoleBinding, ""))
		systemRoleBindings.PUT("/:name", r.metricHandlerWrapper(CreateOrUpdateSystemRoleBinding, ""))
	}

	userBindings := router.Group("userbindings")
	{
		userBindings.GET("", r.metricHandlerWrapper(ListUserBindings, ""))
	}

	bundles := router.Group("bundles")
	{
		bundles.GET("/:name", r.metricHandlerWrapper(DownloadBundle, ""))
	}

	policyRegistrations := router.Group("policies")
	{
		policyRegistrations.PUT("/:resourceName", r.metricHandlerWrapper(CreateOrUpdatePolicyRegistration, ""))
	}

	policyDefinitions := router.Group("policy-definitions")
	{
		policyDefinitions.GET("", r.metricHandlerWrapper(GetPolicyRegistrationDefinitions, ""))
	}

	policySvrHealthz := router.Group("healthz")
	{
		policySvrHealthz.GET("", r.metricHandlerWrapper(Healthz, ""))
	}
	policyUserPermission := router.Group("permission")
	{
		policyUserPermission.GET("/:uid", r.metricHandlerWrapper(GetUserPermission, ""))
	}
}

func (r *Router) metricHandlerWrapper(fn gin.HandlerFunc, apiName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		fn(c)
		r.responseTimeHistogram.WithLabelValues(apiName).Observe(float64(time.Since(start).Milliseconds()))
	}
}
