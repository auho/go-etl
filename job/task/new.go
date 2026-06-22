package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
)

func CleanTask(resource job.CleanResource, modes []mode.UpdateModer, opts ...func(clean *Clean)) {
	clean := NewClean(resource, modes, opts...)
	RunConsumer(resource.Source(), []itemConsumer{clean})
}

func InsertTask(source job.Source, target job.Target, moder mode.InsertModer, opts ...func(*Insert)) {
	insert := NewInsert(target, moder, opts...)
	RunProducer(source, []itemProducer{insert})
}

func TransferTask(source job.Source, target job.Target, moder mode.TransferModer) {
	transfer := NewTransfer(target, moder)
	RunProducer(source, []itemProducer{transfer})
}

func UpdateAndTransferTask(source job.Source, target job.Target, modes []mode.UpdateModer) {
	updateTransfer := NewUpdateAndTransfer(source, target, modes)
	RunProducer(source, []itemProducer{updateTransfer})
}

func UpdateTask(source job.Source, modes []mode.UpdateModer) {
	update := NewUpdate(source, modes)
	RunProducer(source, []itemProducer{update})
}
