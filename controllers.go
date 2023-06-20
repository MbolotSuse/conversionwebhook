package main

import (
	"context"
	"fmt"

	v1 "github.com/Mbolotsuse/conversionwebhook/api/v1"
	v2 "github.com/Mbolotsuse/conversionwebhook/api/v2"
	"github.com/rancher/lasso/pkg/controller"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func runControllers() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("unable to get in cluster config %w", err)
	}

	factory, err := controller.NewSharedControllerFactoryFromConfig(config, scheme)
	if err != nil {
		return fmt.Errorf("unable to construct controller factory: %w", err)
	}
	if err := factory.Start(context.Background(), 3); err != nil {
		return fmt.Errorf("unable to start factory: %w", err)
	}

	hubGVR := schema.GroupVersionResource{Group: "test.cattle.io", Version: "v1", Resource: "foos"}
	hubKind := "Foo"
	hubController := factory.ForResourceKind(hubGVR, hubKind, true)
	hubController.RegisterHandler(context.Background(), "hub-controller", &hubHandler{})
	if err := hubController.Start(context.Background(), 3); err != nil {
		return fmt.Errorf("unable to start controller for hub version %w", err)
	}

	spokeGVR := schema.GroupVersionResource{Group: "test.cattle.io", Version: "v2", Resource: "foos"}
	spokeKind := "Foo"
	spokeController := factory.ForResourceKind(spokeGVR, spokeKind, true)
	spokeController.RegisterHandler(context.Background(), "spoke-controller", &spokeHandler{})
	if err := spokeController.Start(context.Background(), 3); err != nil {
		return fmt.Errorf("unable to start controller for spoke version %w", err)
	}

	klog.Infof("started controllers")
	return nil
}

type hubHandler struct{}

func (h *hubHandler) OnChange(key string, obj runtime.Object) (runtime.Object, error) {
	foo := obj.(*v1.Foo)
	klog.Infof("hub handler Object: %s, initialField: %s, removedField: %s", key, foo.Spec.InitialField, foo.Spec.RemovedField)
	return obj, nil
}

type spokeHandler struct{}

func (h *spokeHandler) OnChange(key string, obj runtime.Object) (runtime.Object, error) {
	foo := obj.(*v2.Foo)
	klog.Infof("spoke handler Object: %s, initialField: %s, addedField: %s", key, foo.Spec.InitialField, foo.Spec.AddedField)
	return obj, nil
}
