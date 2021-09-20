package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gomodules.xyz/jsonpatch/v2"
	"k8s.io/klog"

	v1 "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1"
	admv1 "k8s.io/api/admission/v1"
)

type checkFn func(ar *admv1.AdmissionReview) error
type server struct {
	sv *http.Server
}

var (
	certFile string = "/var/run/secrets/webhook/tls.crt"
	keyFile  string = "/var/run/secrets/webhook/tls.key"
	addr     string = ":9443"
)

func main() {
	flag.StringVar(&certFile, "cert", certFile, "Path to TLS certificate")
	flag.StringVar(&keyFile, "key", keyFile, "Path to TLS certificate key")
	flag.StringVar(&addr, "addr", addr, "Addres to bind the webhook server to")
	flag.Set("logtostderr", "true")
	klog.InitFlags(nil)
	flag.Parse()

	s := newServer()
	if err := s.serve(); err != nil {
		klog.Fatalf("webhook server failed: %v", err)
	}
}

func getPatch(old, updated *admv1.AdmissionReview) ([]byte, *admv1.PatchType, error) {
	klog.V(9).Infof("old raw body: %v", string(old.Request.Object.Raw))
	klog.V(9).Infof("updated raw body: %v", string(updated.Request.Object.Raw))
	patchObj, err := jsonpatch.CreatePatch(old.Request.Object.Raw, updated.Request.Object.Raw)
	if err != nil {
		klog.Errorf("failed to create patch: %v", err)
		return nil, nil, fmt.Errorf("failed to jsonpatch")
	}
	if len(patchObj) == 0 {
		klog.Infof("no patch necessary for request %v/%v", updated.Request.Namespace, updated.Request.Name)
		return nil, nil, nil
	}
	patch, err := json.Marshal(patchObj)
	if err != nil {
		klog.Errorf("failed to marshal JSON patch: %v", err)
		return nil, nil, fmt.Errorf("failed to marshal json")
	}
	pt := admv1.PatchTypeJSONPatch
	klog.Infof("created patch %v, request %v/%v", string(patch), updated.Request.Namespace, updated.Request.Name)
	return patch, &pt, nil
}

func (s *server) mutateClusterCNI(ar *admv1.AdmissionReview) error {
	klog.Infof("setting cluster CNI, request %v/%v", ar.Request.Namespace, ar.Request.Name)

	cluster := &v1.Cluster{}
	if err := json.Unmarshal(ar.Request.Object.Raw, cluster); err != nil {
		klog.Errorf("failed to unmarshal request %v/%v to cluster: %v", ar.Request.Namespace, ar.Request.Name, err)
		return fmt.Errorf("failed to parse request raw object")
	}

	if cluster.Labels["hackaton-cni"] != "" {
		cluster.Spec.CNIPlugin.Type = v1.CNIPluginType(cluster.Labels["hackaton-cni"])
		cluster.Spec.CNIPlugin.Version = ""
	} else {
		klog.Errorf("nothing to do for request %v/%v: %v", ar.Request.Namespace, ar.Request.Name, cluster.Labels)
		return nil
	}

	var err error
	if ar.Request.Object.Raw, err = json.Marshal(cluster); err != nil {
		klog.Errorf("failed to marshal request %v/%v to cluster: %v", ar.Request.Namespace, ar.Request.Name, err)
		return fmt.Errorf("failed to compress request raw object")
	}

	return nil
}

func serve(resp http.ResponseWriter, req *http.Request, check checkFn) {
	klog.Info("serve resp req")
	if req.Body == nil {
		klog.Errorf("req content body is nil")
		return
	}
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	reqReview := &admv1.AdmissionReview{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		klog.Errorf("failed reading req body: %v", err)
		return
	}

	if err := json.Unmarshal(body, reqReview); err != nil {
		klog.Infof("failed parsing request body: %v", string(body))
		klog.Errorf("failed parsing request body: %v", err)
		return
	}

	respReview := &admv1.AdmissionReview{}
	reqReviewCopy := reqReview.DeepCopy()
	klog.Infof("request raw body: %v", string(reqReview.Request.Object.Raw))
	if err := check(reqReview); err == nil {
		patch, patchType, err := getPatch(reqReviewCopy, reqReview)
		if err != nil {
			klog.Errorf("failed making patch: %v", err)
			return
		}
		r := respReview.Response
		r.Patch, r.PatchType = patch, patchType
	}
	respBytes, err := json.Marshal(respReview)
	if err != nil {
		klog.Errorf("failed to marshal response: %v", err)
		return
	}
	if _, err := resp.Write(respBytes); err != nil {
		klog.Errorf("failed to write response: %v", err)
		return
	}
	klog.Infof("finished successfully %v/%v", reqReview.Request.Namespace, reqReview.Request.Name)
}

func (s *server) mutateClusterCNIHandler(resp http.ResponseWriter, req *http.Request) {
	serve(resp, req, s.mutateClusterCNI)
}

func configTLS() *tls.Config {
	sCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		klog.Fatalf("failed to load x509 key pair: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}

func newServer() *server {
	klog.Info("creating cluster CNI hack webhook")
	router := mux.NewRouter()
	s := &server{
		sv: &http.Server{
			Addr:         addr,
			TLSConfig:    configTLS(),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}

	router.HandleFunc("/mutate-cluster-cni", s.mutateClusterCNIHandler)

	s.sv.Handler = router
	return s
}

func (s *server) serve() error {
	return s.sv.ListenAndServeTLS("", "")
}
