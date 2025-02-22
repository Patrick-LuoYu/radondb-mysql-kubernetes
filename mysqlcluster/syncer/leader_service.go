/*
Copyright 2021 RadonDB.

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

package syncer

import (
	"github.com/presslabs/controller-util/syncer"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/radondb/radondb-mysql-kubernetes/mysqlcluster"
	"github.com/radondb/radondb-mysql-kubernetes/utils"
)

// NewLeaderSVCSyncer returns leader service syncer.
func NewLeaderSVCSyncer(cli client.Client, c *mysqlcluster.MysqlCluster) syncer.Interface {
	labels := c.GetLabels()
	labels["mysql.radondb.com/service-type"] = string(utils.LeaderService)
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetNameForResource(utils.LeaderService),
			Namespace: c.Namespace,
			Labels:    labels,
		},
	}
	return syncer.NewObjectSyncer("LeaderSVC", c.Unwrap(), service, cli, func() error {
		// Allows to modify the service access method, the default is ClusterIP.
		if service.Spec.Type == "" {
			service.Spec.Type = "ClusterIP"
		}
		service.Spec.Selector = c.GetSelectorLabels()
		service.Spec.Selector["role"] = string(utils.Leader)

		if len(service.Spec.Ports) != 2 {
			service.Spec.Ports = make([]corev1.ServicePort, 2)
		}

		service.Spec.Ports[0].Name = utils.MysqlPortName
		service.Spec.Ports[0].Port = utils.MysqlPort
		service.Spec.Ports[0].TargetPort = intstr.FromInt(utils.MysqlPort)

		//xtrabckup
		service.Spec.Ports[1].Name = utils.XBackupPortName
		service.Spec.Ports[1].Port = utils.XBackupPort
		service.Spec.Ports[1].TargetPort = intstr.FromInt(utils.XBackupPort)
		return nil
	})
}
