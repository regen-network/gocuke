package gocuke

import "reflect"

func (r *ScenarioContext) StepMethods(suite interface{}) {
	typ := reflect.TypeOf(suite)
	for _, step := range r.pickle.Steps {
		sig := guessMethodSig(step.Text)
		method, ok := typ.MethodByName(sig.name)
		if ok {

		}
	}
}
