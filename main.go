package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type HostEntry struct {
	IP          string
	Hostname    string
	Description string
}

type IPAddressResponse struct {
	Address     string `json:"address"`
	DNSName     string `json:"dns_name"`
	Description string `json:"description"`
}

type IPAddressesListResponse struct {
	Count    int                 `json:"count"`
	Next     string              `json:"next"`
	Previous string              `json:"previous"`
	Results  []IPAddressResponse `json:"results"`
}

var kubeClient *kubernetes.Clientset

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	kubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Error building Kubernetes client: %s", err.Error())
	}

	r := gin.Default()
	r.POST("/webhook", handleWebhook)
	r.Run()
}

func handleWebhook(c *gin.Context) {
	netboxURL := os.Getenv("NETBOX_API_URL")
	netboxToken := os.Getenv("NETBOX_API_TOKEN")

	req, err := http.NewRequest("GET", netboxURL+"/api/ipam/ip-addresses", nil)
	if err != nil {
		logrus.Errorf("Error creating request: %s", err.Error())
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Header.Add("Authorization", "Token "+netboxToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error making request: %s", err.Error())
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Non-OK HTTP response: %s", resp.Status)
		c.JSON(500, gin.H{
			"error": "Non-OK HTTP response: " + resp.Status,
		})
		return
	}

	ipList := &IPAddressesListResponse{}
	if err := json.NewDecoder(resp.Body).Decode(ipList); err != nil {
		logrus.Errorf("Error decoding response: %s", err.Error())
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	hosts := []HostEntry{}
	for _, ip := range ipList.Results {
		hosts = append(hosts, HostEntry{
			IP:          ip.Address,
			Hostname:    ip.DNSName,
			Description: ip.Description,
		})
	}

	hostsFileContent := generateHostsFile(hosts)
	updateConfigMap(hostsFileContent)

	c.JSON(200, gin.H{
		"message": "Hosts file updated successfully",
	})
}

func generateHostsFile(hosts []HostEntry) string {
	entries := []string{}
	for _, host := range hosts {
		entries = append(entries, host.IP+" "+host.Hostname+" # "+host.Description)
	}
	return strings.Join(entries, "\n")
}

func updateConfigMap(content string) {
	namespace := os.Getenv("K8S_NAMESPACE")
	if namespace == "" {
		namespace = "kube-system" // default value
	}
	configMapName := os.Getenv("K8S_CONFIGMAP")
	if configMapName == "" {
		configMapName = "netbox-hosts" // default value
	}

	configMap, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), configMapName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Error getting ConfigMap: %s", err.Error())
		return
	}

	configMap.Data["hosts"] = content
	_, err = kubeClient.CoreV1().ConfigMaps(namespace).Update(context.Background(), configMap, metav1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Error updating ConfigMap: %s", err.Error())
	} else {
		logrus.Info("ConfigMap updated successfully")
	}
}
