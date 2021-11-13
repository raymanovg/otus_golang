package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}

	in = handIn(in, done)
	out := handleStage(stages[0], in, done)

	if len(stages) == 1 {
		return out
	}

	return ExecutePipeline(out, done, stages[1:]...)
}

func handleStage(stage Stage, in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for v := range stage(in) {
			select {
			case <-done:
				return
			case out <- v:
			}
		}
	}()

	return out
}

func handIn(in In, done In) In {
	innerIn := make(Bi)
	go func() {
		defer close(innerIn)
		for v := range in {
			select {
			case <-done:
				return
			case innerIn <- v:
			}
		}
	}()

	return innerIn
}
