package acsengine

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
)

func TestWriteTLSArtifacts(t *testing.T) {

	cs := api.CreateMockContainerService("testcluster", "1.7.12", 1, 2, true)
	writer := &ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: nil,
		},
	}
	dir := "_testoutputdir"
	defaultDir := fmt.Sprintf("%s-%s", cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.GetClusterID())
	defaultDir = path.Join("_output", defaultDir)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(defaultDir)

	// Generate apimodel and azure deploy artifacts without certs
	err := writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", dir, false, false)

	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	expectedFiles := []string{"apimodel.json", "azuredeploy.json", "azuredeploy.parameters.json"}

	for _, f := range expectedFiles {
		if _, err := os.Stat(dir + "/" + f); os.IsNotExist(err) {
			t.Fatalf("expected file %s/%s to be generated by WriteTLSArtifacts", dir, f)
		}
	}

	os.RemoveAll(dir)

	// Generate parameters only and certs
	err = writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", "", true, true)
	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	if _, err := os.Stat(defaultDir + "/apimodel.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/apimodel.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	if _, err := os.Stat(defaultDir + "/azuredeploy.json"); !os.IsNotExist(err) {
		t.Fatalf("expected file %s/azuredeploy.json not to be generated by WriteTLSArtifacts with parametersOnly set to true", defaultDir)
	}

	expectedFiles = []string{"azuredeploy.parameters.json", "ca.crt", "ca.key", "apiserver.crt", "apiserver.key", "client.crt", "client.key", "etcdclient.key", "etcdclient.crt", "etcdserver.crt", "etcdserver.key", "etcdpeer0.crt", "etcdpeer0.key", "kubectlClient.crt", "kubectlClient.key"}

	for _, f := range expectedFiles {
		if _, err := os.Stat(defaultDir + "/" + f); os.IsNotExist(err) {
			t.Fatalf("expected file %s/%s to be generated by WriteTLSArtifacts", dir, f)
		}
	}

	kubeDir := path.Join(defaultDir, "kubeconfig")
	if _, err := os.Stat(kubeDir + "/" + "kubeconfig.eastus.json"); os.IsNotExist(err) {
		t.Fatalf("expected file %s/kubeconfig/kubeconfig.eastus.json to be generated by WriteTLSArtifacts", defaultDir)
	}
	os.RemoveAll(defaultDir)

	// Generate certs with all kubeconfig locations
	cs.Location = ""
	err = writer.WriteTLSArtifacts(cs, "vlabs", "fake template", "fake parameters", "", true, false)
	if err != nil {
		t.Fatalf("unexpected error trying to write TLS artifacts: %s", err.Error())
	}

	for _, region := range helpers.GetAzureLocations() {
		if _, err := os.Stat(kubeDir + "/" + "kubeconfig." + region + ".json"); os.IsNotExist(err) {
			t.Fatalf("expected kubeconfig for region %s to be generated by WriteTLSArtifacts", region)
		}
	}
}
