package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//Cluster type to hold cluster specific details
type Cluster struct {
	restConfig          *rest.Config
	clientset           *kubernetes.Clientset
	apiextensionsClient *apiextensionsclientset.Clientset
}

// FillClientSets checks whether the cluster is accessible and returns an error if not
func (c *Cluster) FillClientSets() error {
	clientset, err := kubernetes.NewForConfig(c.restConfig)
	if err == nil {
		c.clientset = clientset
		err = c.ProbeNodes()
		if err != nil {
			return err
		}
		apiextensionsClient, err := apiextensionsclientset.NewForConfig(c.restConfig)
		if err == nil {
			c.apiextensionsClient = apiextensionsClient
		}
	}
	return err
}

//ProbeNodes tries to probe for nodes. Indirectly it checks whether the cluster is accessible.
// If not accessible, then it returns an error
func (c *Cluster) ProbeNodes() error {
	_, err := c.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	return err
}

// NewCluster returns a Cluster struct
func NewCluster(kubeConfigPath string) (c *Cluster, e error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err == nil {
		c = &Cluster{
			restConfig: config,
		}
	} else {
		c = &Cluster{}
	}

	return c, err
}

// GetClientset returns a Clientset
func (c *Cluster) GetClientset() (k *kubernetes.Clientset) {
	return c.clientset
}

func parseK8sYaml(filepath string) ([]runtime.Object, []*schema.GroupVersionKind, error) {
	fileR, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, nil, err
	}
	acceptedK8sTypes := regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount|CustomResourceDefinition)`)
	fileAsString := string(fileR[:])
	sepYamlfiles := strings.Split(fileAsString, "---")
	retObj := make([]runtime.Object, 0, len(sepYamlfiles))
	retKind := make([]*schema.GroupVersionKind, 0, len(sepYamlfiles))
	for _, f := range sepYamlfiles {
		if f == "\n" || f == "" {
			// ignore empty cases
			continue
		}

		decode := scheme.Codecs.UniversalDeserializer().Decode
		obj, groupVersionKind, err := decode([]byte(f), nil, nil)

		if err != nil {
			log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
			continue
		}

		if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
			log.Printf("The custom-roles configMap contained K8s object types which are not supported! Skipping object with type: %s", groupVersionKind.Kind)
		} else {
			retKind = append(retKind, groupVersionKind)
			retObj = append(retObj, obj)
		}

	}
	return retObj, retKind, err
}

// ApplyYamlFile uses yaml to create resources in kubernetes
func (c *Cluster) ApplyYamlFile(filePath string) error {
	/* TO-DO: This function checks for the availability of filePath
	if available, then uses kubectl to perform kubectl apply on that file
	*/
	runtimeobj, kind, err := parseK8sYaml(filePath)
	if err != nil {
		for key, obj := range runtimeobj {
			switch kind[key].Kind {
			case "CustomResourceDefinition":
				crd := obj.(*v1beta1.CustomResourceDefinition)
				_, err := c.apiextensionsClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
