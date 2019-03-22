package automulticlusterapp

import (
	"context"
	"fmt"

	rancheroperatorv1alpha1 "gitlab.thalesdigital.io/core-kube/rancher-operator/pkg/apis/rancheroperator/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	managementrancherv3 "github.com/rancher/types/apis/management.cattle.io/v3"
)

var log = logf.Log.WithName("controller_automulticlusterapp")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AutoMultiClusterApp Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) *ReconcileAutoMultiClusterApp {
	return &ReconcileAutoMultiClusterApp{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileAutoMultiClusterApp) error {
	// Create a new controller
	c, err := controller.New("automulticlusterapp-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// Watch for changes to primary resource AutoMultiClusterApp
	// This resources is an APIS object for this controller created with operator sdk
	err = c.Watch(&source.Kind{Type: &rancheroperatorv1alpha1.AutoMultiClusterApp{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileAutoMultiClusterApp{}

// ReconcileAutoMultiClusterApp reconciles a AutoMultiClusterApp object
type ReconcileAutoMultiClusterApp struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AutoMultiClusterApp object and makes changes based on the state read
// and what is in the AutoMultiClusterApp.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAutoMultiClusterApp) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Name", request.Name)

	ctx := context.TODO()

	// Our resource is cluster scoped
	request.NamespacedName.Namespace = ""

	reqLogger.Info("Reconciling AutoMultiClusterApp")

	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// Fetch the AutoMultiClusterApp instance
	// !!!! Here we have a result not empty even if it's not a core resources !!!!
	instance := &rancheroperatorv1alpha1.AutoMultiClusterApp{}
	err := r.client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// Now let's try to list items for a Third party Resources (From Rancher)
	// https://github.com/rancher/types/tree/master/apis/management.cattle.io/v3
	// We add the third party resource by following the documentation :
	// https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#adding-3rd-party-resources-to-your-operator

	projects := &managementrancherv3.ProjectList{}
	opt := &client.ListOptions{
		Raw: &metav1.ListOptions{
			ResourceVersion: "0",
		},
	}
	opt.InNamespace("")

	if err := r.client.List(ctx, opt, projects); err != nil {
		reqLogger.Info("Failed to list projects")
		return reconcile.Result{}, err
	}
	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// With default client list is always empty (ONLY when we use a third party resource)
	// With Core resources (e.g POD) it's working.
	fmt.Println("---First Try using default client---")
	reqLogger.Info("List projects", "project", len(projects.Items))

	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// Now let's try to create another client by following documentatiion :
	//https://github.com/operator-framework/operator-sdk/blob/master/doc/user/client.md#non-default-client
	cfg, _ := config.GetConfig()
	c, _ := client.New(cfg, client.Options{
		Scheme: r.scheme,
	})

	projects2 := &managementrancherv3.ProjectList{}
	opts := &client.ListOptions{}

	if err := c.List(ctx, opts, projects2); err != nil {
		reqLogger.Info("Failed to list projects")
		return reconcile.Result{}, err
	}

	// !!! NOTE FOR OPERATOR SDK TEAMS !!!
	// Result is empty but we have the resources present on the cluster.
	fmt.Println("---SECOND Try using default client---")
	reqLogger.Info("List projects", "projects", len(projects2.Items))

	return reconcile.Result{}, nil
}
