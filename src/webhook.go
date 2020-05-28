package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "time"

    "encoding/json"
    v1beta1 "k8s.io/api/admission/v1beta1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Mutate mutates
func mutateTimezone(body []byte, verbose bool) ([]byte, error) {

	// unmarshal request into AdmissionReview struct
	admReview := v1beta1.AdmissionReview{}
	if err := json.Unmarshal(body, &admReview); err != nil {
		return nil, fmt.Errorf("unmarshaling request failed with %s", err)
	}

	var err error
	var pod *corev1.Pod

	responseBody := []byte{}
	ar := admReview.Request
	resp := v1beta1.AdmissionResponse{}

	if ar != nil {

		// get the Pod object and unmarshal it into its struct, if we cannot, we might as well stop here
		if err := json.Unmarshal(ar.Object.Raw, &pod); err != nil {
			return nil, fmt.Errorf("unable unmarshal pod json object %v", err)
		}
		// set response options
		resp.Allowed = true
		resp.UID = ar.UID
		pT := v1beta1.PatchTypeJSONPatch
		resp.PatchType = &pT // it's annoying that this needs to be a pointer as you cannot give a pointer to a constant?

		// add some audit annotations, helpful to know why a object was modified, maybe (?)
		resp.AuditAnnotations = map[string]string{
			"timezone-webhook": "modified",
		}

		// the actual mutation is done by a string in JSONPatch style, i.e. we don't _actually_ modify the object, but
		// tell K8S how it should modifiy it
		p := []map[string]string{}

		volumeMount := []corev1.VolumeMount{"timezone","true","/etc/localtime"}

        for i := range pod.Spec.Containers {
            patch := patchOperation{
                Op:    "add",
                Path:  fmt.Sprintf("/spec/containers/%d/volumeMounts", i),
                Value: volumeMount,
            }
            p = append(p, patch)
        }

        volumeSource := []corev1.VolumeSource{"/etc/localtime"}
        volume := []corev1.Volume{"timezone",volumeSource}

        patchVolumes := patchOperation{
            Op:    "add",
            Path:  "/spec/volumes",
            Value: volume,
        }
        p = append(p, patchVolumes)

		// parse the []map into JSON
		resp.Patch, err = json.Marshal(p)

		resp.Result = &metav1.Status{
			Status: "Success",
		}

		admReview.Response = &resp
		// back into JSON so we can return the finished AdmissionReview w/ Response directly
		// w/o needing to convert things in the http handler
		responseBody, err = json.Marshal(admReview)
		if err != nil {
			return nil, err // untested section
		}
	}

	if verbose {
		log.Printf("resp: %s\n", string(responseBody)) // untested section
	}

	return responseBody, nil
}

func handleMutateTimezone(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		sendError(err, w)
		return
	}

	mutated, err := mutateTimezone(body, true)
	if err != nil {
		sendError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(mutated)
}

func sendError(err error, w http.ResponseWriter) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", err)
}

func main() {
    log.Println("Starting webhook server ...")

    mux := http.NewServeMux()
    mux.HandleFunc("/mutateTimezone", handleMutateTimezone)

    s := &http.Server{
        Addr:           "8080",
        Handler:        mux,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20, // 1048576
    }

    s.ListenAndServe()
}