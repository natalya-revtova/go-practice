package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	emptyOut := func() Out {
		out := make(Bi)
		close(out)
		return out
	}

	if in == nil {
		return emptyOut()
	}
	if stages == nil {
		return emptyOut()
	}

	for _, stage := range stages {
		in = runner(done, stage(in))
	}
	return in
}

func runner(done In, in In) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for value := range in {
			select {
			case <-done:
				return
			default:
				out <- value
			}
		}
	}()

	return out
}
