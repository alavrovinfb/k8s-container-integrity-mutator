package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScienceSoft-Inc/k8s-container-integrity-mutator/pkg/mutate"
	"github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

const MimeTypeApplicationJson = "application/json"

type Router interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

type Handers struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *Handers {
	return &Handers{logger: logger}
}

func (h *Handers) Register(r Router) {
	r.HandleFunc("/mutate", h.mutate)
}

func (h *Handers) mutate(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("starting mutation")
	admReview, err := admissionReviewFromRequest(r)
	if err != nil {
		h.logger.WithError(err).Errorf("Failed getting admission review from request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	admResp, err := mutate.InjectIntegrityMonitor(h.logger, admReview.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// the final response will be another admission review
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admResp
	admissionReviewResponse.Response.UID = admReview.Request.UID
	admissionReviewResponse.SetGroupVersionKind(admReview.GroupVersionKind())

	resp, err := json.Marshal(admissionReviewResponse)
	if err != nil {
		h.logger.WithError(err).Errorf("Failed marshal admission review from response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", MimeTypeApplicationJson)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func admissionReviewFromRequest(r *http.Request) (*admissionv1.AdmissionReview, error) {
	// Validate that the incoming content type is correct.
	if r.Header.Get("Content-Type") != MimeTypeApplicationJson {
		return nil, fmt.Errorf("expected application/json content-type")
	}

	var admissionReviewRequest admissionv1.AdmissionReview
	err := json.NewDecoder(r.Body).Decode(&admissionReviewRequest)
	if err != nil {
		return nil, err
	}
	return &admissionReviewRequest, nil
}
