package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	fhir "github.com/intervention-engine/fhir/models"
)

func UploadResources(resources []interface{}, baseURL string) (map[string]string, error) {
	refMap := make(map[string]string)
	for _, t := range resources {
		updateReferences(t, refMap)
		newLoc, err := UploadResource(t, baseURL)
		if err != nil {
			return refMap, err
		}

		// Add entry to map from oldid to http://url/to/new/resource
		id := reflect.ValueOf(t).FieldByName("Id").String()
		if id != "" {
			refMap[id] = newLoc
		}
	}

	return refMap, nil
}

func UploadResource(resource interface{}, baseURL string) (string, error) {
	url := baseURL + "/" + reflect.TypeOf(resource).Name()
	json, _ := json.Marshal(resource)
	body := bytes.NewReader(json)
	response, err := http.Post(url, "application/json+fhir", body)
	if err != nil {
		return "", err
	}

	return response.Header.Get("Location"), nil
}

func updateReferences(resource interface{}, refMap map[string]string) {
	s := reflect.ValueOf(resource)
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if f.Type() == reflect.TypeOf(&fhir.Reference{}) && !f.IsNil() {
			updateReference(f.Interface().(*fhir.Reference), refMap)
		} else if f.Type() == reflect.TypeOf([]fhir.Reference{}) {
			for i := 0; i < f.Len(); i++ {
				updateReference(f.Index(i).Addr().Interface().(*fhir.Reference), refMap)
			}
		}
	}
}

func updateReference(ref *fhir.Reference, refMap map[string]string) {
	if ref != nil {
		newRef, ok := refMap[strings.TrimPrefix(ref.Reference, "cid:")]
		if ok {
			ref.Reference = newRef
		} else {
			panic(fmt.Sprint("Failed to find updated reference for ", ref))

		}
	}
}
