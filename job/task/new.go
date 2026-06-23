package task

import (
	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/transform"
)

func CleanTask(resource job.CleanResource, modes []transform.UpdateOperator, opts ...func(clean *Clean)) {
	clean := NewClean(resource, modes, opts...)
	RunConsumer(resource.Source(), []itemConsumer{clean})
}

func InsertTask(source job.Source, target job.Target, moder transform.InsertOperator, opts ...func(*Insert)) {
	insert := NewInsert(target, moder, opts...)
	RunProducer(source, []itemProducer{insert})
}

func TransferTask(source job.Source, target job.Target, moder transform.TransferOperator) {
	transfer := NewTransfer(target, moder)
	RunProducer(source, []itemProducer{transfer})
}

func UpdateAndTransferTask(source job.Source, target job.Target, modes []transform.UpdateOperator) {
	updateTransfer := NewUpdateAndTransfer(source, target, modes)
	RunProducer(source, []itemProducer{updateTransfer})
}

func UpdateTask(source job.Source, modes []transform.UpdateOperator) {
	update := NewUpdate(source, modes)
	RunProducer(source, []itemProducer{update})
}
