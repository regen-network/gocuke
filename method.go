package gocuke

import "reflect"

func (r *ScenarioContext) RegisterStepMethods(suite interface{}) {
	typ := reflect.TypeOf(suite)
	for _, step := range r.pickle.Steps {
		sig := guessMethodSig(step.Text)
		method, ok := typ.MethodByName(sig.name)
		if ok {
			numIn := method.Type.NumIn()
			in := make([]reflect.Type, numIn-1)
			for i := 1; i < numIn; i++ {
				in[i-1] = method.Type.In(i)
			}

			newTy := reflect.FuncOf(in, nil, false)
			r.stepDefs = append(r.stepDefs, &stepDef{
				exp: sig.regex,
				f: reflect.MakeFunc(newTy, func(args []reflect.Value) (results []reflect.Value) {
					newArgs := []reflect.Value{reflect.ValueOf(suite)}
					newArgs = append(newArgs, args...)
					return method.Func.Call(newArgs)
				}),
			})
		}
	}
}
