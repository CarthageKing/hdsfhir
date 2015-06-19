package hdsfhir

import fhir "github.com/intervention-engine/fhir/models"

type Condition struct {
	Entry
	Severity CodeMap `json:"severity"`
}

func (c *Condition) FHIRModels() []interface{} {
	fhirCondition := &fhir.Condition{Id: c.GetTempID()}
	fhirCondition.Code = c.Codes.FHIRCodeableConcept(c.Description)
	fhirCondition.OnsetDateTime = c.StartTime.FHIRDateTime()
	fhirCondition.Patient = c.Patient.FHIRReference()
	if c.EndTime != 0 {
		fhirCondition.AbatementDate = c.EndTime.FHIRDateTime()
	}

	return []interface{}{fhirCondition}
}
