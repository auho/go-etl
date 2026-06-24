package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
)

func CleanTask(resource job.CleanResource, modes []transform.UpdateOperator, opts ...func(clean *Clean)) {
	clean := NewClean(resource, modes, opts...)
	RunConsumer(resource.Source(), []itemConsumer{clean})
}

func InsertTask(source job.Table, target job.Table, moder transform.InsertOperator, opts ...func(*Insert)) {
	insert := NewInsert(target, moder, opts...)
	RunProducer(source, []itemProducer{insert})
}

func TransferTask(source job.Table, target job.Table, moder transform.TransferOperator) {
	transfer := NewTransfer(target, moder)
	RunProducer(source, []itemProducer{transfer})
}

func UpdateTransferTask(source job.Table, target job.Table, modes []transform.UpdateOperator) {
	updateTransfer := NewUpdateTransfer(source, target, modes)
	RunProducer(source, []itemProducer{updateTransfer})
}

func UpdateTask(source job.Table, modes []transform.UpdateOperator) {
	update := NewUpdate(source, modes)
	RunProducer(source, []itemProducer{update})
}
