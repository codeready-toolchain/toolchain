package toolchainclusterresources

import (
	"context"
	"embed"
	"fmt"

	toolchainv1alpha1 "github.com/codeready-toolchain/api/api/v1alpha1"
	commoncontroller "github.com/codeready-toolchain/toolchain-common/controllers"
	applycl "github.com/codeready-toolchain/toolchain-common/pkg/client"
	commonpredicates "github.com/codeready-toolchain/toolchain-common/pkg/predicate"
	"github.com/codeready-toolchain/toolchain-common/pkg/template"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ResourceControllerLabelValue is being added to all the resources managed by this controller.
// It's then used to filter all the events on those resources by using a mapper function in the watcher configuration.
const ResourceControllerLabelValue = "toolchaincluster-resources-controller" // TODO move this label value to api repo

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager, operatorNamespace string) error {
	// check for required templates FS directory
	if r.Templates == nil {
		return fmt.Errorf("no templates FS configured")
	}

	build := ctrl.NewControllerManagedBy(mgr).
		For(&v1.ServiceAccount{})

	// add watcher for all kinds from given templates
	var err error
	r.templateObjects, err = template.LoadObjectsFromEmbedFS(r.Templates, &template.Variables{Namespace: operatorNamespace})
	if err != nil {
		return err
	}
	mapToOwnerByLabel := handler.EnqueueRequestsFromMapFunc(commoncontroller.MapToControllerByMatchingLabel(toolchainv1alpha1.ProviderLabelKey, ResourceControllerLabelValue))
	for _, obj := range r.templateObjects {
		build = build.Watches(obj.DeepCopyObject().(runtimeclient.Object), mapToOwnerByLabel, builder.WithPredicates(commonpredicates.LabelsAndGenerationPredicate{}))
	}
	return build.Complete(r)
}

// Reconciler reconciles a ToolchainCluster object
type Reconciler struct {
	Client          runtimeclient.Client
	Scheme          *runtime.Scheme
	Templates       *embed.FS
	templateObjects []*unstructured.Unstructured
}

// Reconcile loads all the manifests from a given embed.FS folder, evaluates the supported variables and applies the objects in the cluster.
func (r *Reconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx)
	reqLogger.Info("Reconciling ToolchainCluster resources controller")
	// check for required templates FS directory
	if r.Templates == nil {
		return reconcile.Result{}, fmt.Errorf("no templates FS configured")
	}

	// apply all the objects with a custom label
	newLabels := map[string]string{
		toolchainv1alpha1.ProviderLabelKey: ResourceControllerLabelValue,
	}

	// TODO implement delete logic for objects that were renamed/removed from the templates

	return reconcile.Result{}, applycl.ApplyUnstructuredObjectsWithNewLabels(ctx, r.Client, r.templateObjects, newLabels) // apply objects on the cluster
}
