package main

import (
	"context"
	"fmt"
	kidlev1beta1 "github.com/orphaner/kidle/pkg/api/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KidleClient struct {
	client.Client
	Namespace string
}

func NewKidleClient(namespaceOverride string) (*KidleClient, error) {
	var restConfig *rest.Config
	var err error

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	configOverrides := &clientcmd.ConfigOverrides{
		Context: clientcmdapi.Context{
			Namespace: namespaceOverride,
		},
	}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return nil, fmt.Errorf("error when getting current namespace: %v", err)
	}

	restConfig, err = kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("error when creating restConfig: %v", err)
	}

	// Create a client with kidle scheme registered
	kidlev1beta1.AddToScheme(scheme.Scheme)
	client, err := client.New(restConfig, client.Options{})
	return &KidleClient{Client: client, Namespace: namespace}, nil
}

func (k *KidleClient) applyDesiredIdleState(idle bool, req *client.ObjectKey) (bool, error) {

	ctx := context.Background()

	// get the IdlingResource from the req
	ir := kidlev1beta1.IdlingResource{}
	err := k.Get(ctx, *req, &ir)
	if err != nil {
		return false, fmt.Errorf("unable to get idlingresource: %v", err)
	}

	// nothing to do if current state == desired state
	if ir.Spec.Idle == idle {
		return false, nil
	}

	// update idle flag to desired state
	ir.Spec.Idle = idle

	err = k.Update(ctx, &ir)
	if err != nil {
		return false, fmt.Errorf("unable to update idlingresource: %v", err)
	}
	return true, nil
}
