package configuration

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewaycontroller "github.com/kong/kubernetes-ingress-controller/v3/internal/controllers/gateway"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/gatewayapi"
	kongv1beta1 "github.com/kong/kubernetes-ingress-controller/v3/pkg/apis/configuration/v1beta1"
)

const maxNAncestors = 16

type serviceStatus struct {
	service           corev1.Service
	acceptedCondition metav1.Condition
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Reconciler Helpers
// -----------------------------------------------------------------------------

// enforceKongUpstreamPolicyStatus gets a list of services (ancestors) along with their desired status and enforce them
// in the KongUpstreamPolicy status.
func (r *KongUpstreamPolicyReconciler) enforceKongUpstreamPolicyStatus(ctx context.Context, oldPolicy *kongv1beta1.KongUpstreamPolicy) (bool, error) {
	// get all the services that reference this UpstreamPolicy
	services := &corev1.ServiceList{}
	err := r.List(ctx, services,
		client.InNamespace(oldPolicy.Namespace),
		client.MatchingFields{
			upstreamPolicyIndexKey: oldPolicy.Name,
		},
	)
	if err != nil {
		return false, err
	}

	// build the desired KongUpstreamPolicy status
	servicesStatus, err := r.buildServicesStatus(ctx, k8stypes.NamespacedName{
		Namespace: oldPolicy.Namespace,
		Name:      oldPolicy.Name,
	}, services.Items)
	if err != nil {
		return false, err
	}

	newPolicyStatus := gatewayapi.PolicyStatus{}
	if len(servicesStatus) > 0 {
		newPolicyStatus.Ancestors = make([]gatewayapi.PolicyAncestorStatus, 0, len(servicesStatus))
	}
	for _, ss := range servicesStatus {
		newPolicyStatus.Ancestors = append(newPolicyStatus.Ancestors,
			gatewayapi.PolicyAncestorStatus{
				AncestorRef: gatewayapi.ParentReference{
					Group:     lo.ToPtr(gatewayapi.Group("core")),
					Kind:      lo.ToPtr(gatewayapi.Kind("Service")),
					Namespace: lo.ToPtr(gatewayapi.Namespace(ss.service.Namespace)),
					Name:      gatewayapi.ObjectName(ss.service.Name),
				},
				ControllerName: gatewaycontroller.GetControllerName(),
				Conditions: []metav1.Condition{
					ss.acceptedCondition,
				},
			},
		)
	}
	if isStatusUpdated := isPolicyStatusUpdated(oldPolicy.Status, newPolicyStatus); !isStatusUpdated {
		newPolicy := oldPolicy.DeepCopy()
		newPolicy.Status = newPolicyStatus
		return true, r.Client.Status().Patch(ctx, newPolicy, client.MergeFrom(oldPolicy))
	}
	return false, nil
}

// buildServicesStatus creates a list of services with their conditions associated.
func (r *KongUpstreamPolicyReconciler) buildServicesStatus(ctx context.Context, upstreamPolicyNN k8stypes.NamespacedName, services []corev1.Service) ([]serviceStatus, error) {
	// sort the services by creationTimestamp
	sort.Slice(services, func(i, j int) bool {
		return services[i].CreationTimestamp.Before(&services[j].CreationTimestamp)
	})

	type indexedServiceStatus struct {
		index int
		data  serviceStatus
	}

	// prepare a service mapping to be used in subsequent operations
	mappedServices := make(map[string]indexedServiceStatus)
	for i, service := range services {
		if i < maxNAncestors {
			acceptedCondition := metav1.Condition{
				Type:               string(gatewayapi.PolicyConditionAccepted),
				Status:             metav1.ConditionTrue,
				Reason:             string(gatewayapi.PolicyReasonAccepted),
				LastTransitionTime: metav1.Now(),
			}
			mappedServices[buildServiceReference(service.Namespace, service.Name)] = indexedServiceStatus{
				index: i,
				data: serviceStatus{
					service:           service,
					acceptedCondition: acceptedCondition,
				},
			}
		} else {
			r.Log.Info(fmt.Sprintf("kongUpstreamPolicy %s/%s status has already %d ancestors, cannot set service %s/%s as an ancestor in the status",
				upstreamPolicyNN.Namespace,
				upstreamPolicyNN.Name,
				maxNAncestors,
				service.Namespace,
				service.Name))
		}
	}

	for serviceKey, serviceStatus := range mappedServices {
		httpRoutes := &gatewayapi.HTTPRouteList{}
		err := r.List(ctx, httpRoutes,
			client.MatchingFields{
				routeBackendRefServiceNameIndexKey: serviceKey,
			},
		)
		if err != nil {
			return nil, err
		}

		for _, httpRoute := range httpRoutes.Items {
			for _, rule := range httpRoute.Spec.Rules {
				if len(rule.BackendRefs) == 0 {
					continue
				}
				for _, br := range rule.BackendRefs {
					serviceRef := backendRefToServiceRef(httpRoute.Namespace, br.BackendRef)
					if serviceRef == "" {
						continue
					}
					if _, ok := mappedServices[serviceRef]; !ok {
						serviceStatus.data.acceptedCondition.Status = metav1.ConditionFalse
						serviceStatus.data.acceptedCondition.Reason = string(gatewayapi.PolicyReasonConflicted)
						mappedServices[serviceKey] = serviceStatus
					}
				}
			}
		}
	}

	servicesStatus := make([]serviceStatus, len(mappedServices))
	for _, ms := range mappedServices {
		servicesStatus[ms.index] = ms.data
	}
	return servicesStatus, nil
}

// -----------------------------------------------------------------------------
// KongUpstreamPolicy Controller - Helpers
// -----------------------------------------------------------------------------

func isPolicyStatusUpdated(oldStatus, newStatus gatewayapi.PolicyStatus) bool {
	if len(oldStatus.Ancestors) != len(newStatus.Ancestors) {
		return false
	}
	for i, oldAncestor := range oldStatus.Ancestors {
		newAncestor := newStatus.Ancestors[i]
		if newAncestor.ControllerName != oldAncestor.ControllerName {
			return false
		}
		if !reflect.DeepEqual(newAncestor.AncestorRef, oldAncestor.AncestorRef) {
			return false
		}

		if len(oldAncestor.Conditions) != len(newAncestor.Conditions) {
			return false
		}
		for j, oldCondition := range oldAncestor.Conditions {
			newCondition := newAncestor.Conditions[j]
			if newCondition.Type != oldCondition.Type ||
				newCondition.Status != oldCondition.Status ||
				newCondition.Reason != oldCondition.Reason ||
				newCondition.Message != oldCondition.Message {
				return false
			}
		}
	}

	return true
}

func backendRefToServiceRef(routeNamespace string, br gatewayapi.BackendRef) string {
	if br.Group != nil && *br.Group != "" && *br.Group != "core" {
		return ""
	}
	if br.Kind != nil && *br.Kind != "" && *br.Kind != "Service" {
		return ""
	}
	namespace := routeNamespace
	if br.Namespace != nil {
		namespace = string(*br.Namespace)
	}
	return buildServiceReference(namespace, string(br.Name))
}

func buildServiceReference(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}