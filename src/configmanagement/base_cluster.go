package configmanagement

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s-config-mgmt/src/config"
	"fmt"
)

type K8sConfigManagementBaseCluster struct {
	K8sConfigManagementBase
	funcs K8sConfigManagementBaseClusterFuncs
}

type K8sConfigManagementBaseClusterFuncs interface {
	init()
	listExistingItems() (map[string]runtime.Object, error)
	listConfigItems() (map[string]config.ConfigObject)
	deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object)
	handleCreate(k8sItem runtime.Object) (error)
	handleUpdate(k8sItem runtime.Object) (error)
	handleDelete(k8sItem runtime.Object) (error)
}

func (mgmt *K8sConfigManagementBaseCluster) Manage() {
	mgmt.funcs.init()

	configList := mgmt.funcs.listConfigItems()

	// check if anything is to do
	if !mgmt.Configuration.AutoCleanup && len(configList) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.funcs.listExistingItems()
	if err != nil {
		mgmt.Logger.StepResult(fmt.Sprintf("[ERROR] %v", err))
		return
	}

	for _, item := range configList {
		if mgmt.isK8sObjectFiltered(item.Object) {
			mgmt.Logger.Step("ignoring %v (filtered)", item.Name)
			continue
		}

		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("updating %v", item.Name)

			// update
			updatedObject := mgmt.funcs.deepCloneObject(item.Object, k8sObject)

			if mgmt.IsNotDryRun() {
				err := mgmt.funcs.handleUpdate(*updatedObject)

				if (mgmt.Force && err != nil) {
					mgmt.Logger.StepResult("failed, forcing recreate")
					mgmt.handleOperationState(mgmt.funcs.handleDelete(*updatedObject))
					mgmt.handleOperationState(mgmt.funcs.handleCreate(item.Object))
				} else {
					mgmt.handleOperationState(err)
				}
			}

		} else {
			mgmt.Logger.Step("creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				mgmt.handleOperationState(mgmt.funcs.handleCreate(item.Object))
			}
		}
	}

	// cleanup
	if mgmt.Configuration.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := configList[k8sObject.(v1.Object).GetName()]; !ok {
				if ! mgmt.isK8sObjectFiltered(k8sObject) {
					mgmt.Logger.Step("deleting %v", k8sObject.(v1.Object).GetName())

					if mgmt.IsNotDryRun() {
						mgmt.handleOperationState(mgmt.funcs.handleDelete(k8sObject))
					}
				} else {
					mgmt.Logger.Step("ignoring %v (filtered)", k8sObject.(v1.Object).GetName())
				}
			}
		}
	}
}
