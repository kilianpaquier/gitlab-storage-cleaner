package artifacts

import (
	"github.com/fogfactory/pipe"
	"github.com/samber/lo"
)

// PipeProjectBuilder is the pipe builder for Project structure.
type PipeProjectBuilder struct {
	preproc    []pipe.Process[Project]
	dispatched *pipe.PoolProcess[Project]
	postproc   []pipe.Process[Project]
}

// NewPipeProjectBuilder creates a new PipeProjectBuilder.
func NewPipeProjectBuilder() *PipeProjectBuilder {
	return &PipeProjectBuilder{}
}

// Processor adds a processor to either the pre processors or post processors of Project structure.
func (b *PipeProjectBuilder) Processor(proc pipe.Process[Project]) *PipeProjectBuilder {
	if b.dispatched == nil {
		b.preproc = append(b.preproc, proc)
	} else {
		b.postproc = append(b.postproc, proc)
	}
	return b
}

// Split defines the split function to send Projects' jobs into pipe.
func (b *PipeProjectBuilder) Split(split pipe.Split[Project, Job]) *PipeJobBuilder {
	return &PipeJobBuilder{parent: b, split: split}
}

// Build builds the pipe processors of Project structure.
func (b *PipeProjectBuilder) Build() pipe.PoolProcess[Project] {
	if b.dispatched == nil {
		panic("no dispatcher")
	}

	preprocs := pipe.Link(pipe.AsPoolProcesses(b.preproc...)...)
	postprocs := pipe.Link(pipe.AsPoolProcesses(b.postproc...)...)
	return pipe.Link(preprocs, *b.dispatched, postprocs)
}

// PipeJobBuilder is the builder for Job pipe processing.
type PipeJobBuilder struct {
	parent *PipeProjectBuilder

	split pipe.Split[Project, Job]
	procs []pipe.Process[Job]
}

// Processor adds a processor to Job structure.
func (b *PipeJobBuilder) Processor(proc pipe.Process[Job]) *PipeJobBuilder {
	b.procs = append(b.procs, proc)
	return b
}

// Merge defines the merge function from Jobs' to their project.
// It returns the parent PipeProjectBuilder.
func (b *PipeJobBuilder) Merge(merge pipe.Merge[Project, Job]) *PipeProjectBuilder {
	dispatcher, err := pipe.NewDispatch(b.split, merge)
	if err != nil {
		panic(err)
	}

	procs := pipe.Link(pipe.AsPoolProcesses(b.procs...)...)
	b.parent.dispatched = lo.ToPtr(pipe.Wrap(procs, dispatcher))
	return b.parent
}
