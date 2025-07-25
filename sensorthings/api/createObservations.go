package api

import (
	entities "github.com/FObersteiner/gosta-core"
)

// PostCreateObservations checks for correctness of the datastreams and observations and calls PostcreateObservations on the database
// ToDo: use transactions
func (a *APIv1) PostCreateObservations(data *entities.CreateObservations) ([]string, []error) {
	_, err := containsMandatoryParams(data)
	if err != nil {
		return nil, err
	}

	returnList := make([]string, 0)

	for i := range len(data.Datastreams) {
		for j := range len(data.Datastreams[i].Observations) {
			obs, errors := a.PostObservationByDatastream(data.Datastreams[i].ID, data.Datastreams[i].Observations[j])
			if errors == nil || len(errors) == 0 {
				returnList = append(returnList, obs.GetSelfLink())
			} else {
				errorString := ""
				for k := range len(errors) {
					if len(errorString) > 0 {
						errorString += ", "
					}

					errorString += errors[k].Error()
				}

				returnList = append(returnList, errorString)
			}
		}
	}

	return returnList, nil
}
