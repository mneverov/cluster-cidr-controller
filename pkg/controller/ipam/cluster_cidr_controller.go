package ipam

import (
	"context"

	v1 "github.com/mneverov/cluster-cidr-controller/pkg/apis/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=networking.x-k8s.io,resources=clustercidrs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

type ClusterCIDRReconciler struct {
	client                  client.Client
	recorder                record.EventRecorder
	multiCIDRRangeAllocator multiCIDRRangeAllocator
}

type ClusterCIDRReconcilerOptions struct {
	Client                  client.Client
	MultiCIDRRangeAllocator multiCIDRRangeAllocator
}

func NewClusterCidrReconciler(opts ClusterCIDRReconcilerOptions) *ClusterCIDRReconciler {
	return &ClusterCIDRReconciler{
		client:                  opts.Client,
		multiCIDRRangeAllocator: opts.MultiCIDRRangeAllocator,
	}
}

func (r *ClusterCIDRReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cidr := &v1.ClusterCIDR{}
	if err := r.client.Get(ctx, req.NamespacedName, cidr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !cidr.GetDeletionTimestamp().IsZero() {
		return r.delete(ctx, cidr)
	}

	return r.reconcile(ctx, cidr)
}

func (r *ClusterCIDRReconciler) reconcile(ctx context.Context, cidr *v1.ClusterCIDR) (ctrl.Result, error) {
	err := r.multiCIDRRangeAllocator.syncClusterCIDR(ctx, cidr)
	return ctrl.Result{}, err
}

func (r *ClusterCIDRReconciler) delete(ctx context.Context, cidr *v1.ClusterCIDR) (ctrl.Result, error) {
	err := r.multiCIDRRangeAllocator.reconcileDelete(ctx, cidr)
	return ctrl.Result{}, err
}

func (r *ClusterCIDRReconciler) enqueueRequestsFromNode(ctx context.Context, object client.Object) (requests []reconcile.Request) {
	node, ok := object.(*corev1.Node)
	if !ok {
		return requests
	}
	// skip if node is null or does not have PodCIDRs
	if node == nil || len(node.Spec.PodCIDRs) == 0 {
		return nil
	}
	cidr, err := r.multiCIDRRangeAllocator.allocatedClusterCIDR(node)
	if err != nil {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{Name: cidr.Name},
	}}
}

func (r *ClusterCIDRReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("clustercidr-reconciler")

	return ctrl.NewControllerManagedBy(mgr).
		Named("clustercidr-reconciler").
		For(&v1.ClusterCIDR{}).
		Watches(&corev1.Node{}, handler.EnqueueRequestsFromMapFunc(r.enqueueRequestsFromNode)).
		Complete(r)
}
