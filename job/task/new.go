package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
)

func CleanTask(resource job.CleanResource, modes []transform.UpdateOperator, opts ...func(clean *Clean)) error {
	clean := NewClean(resource, modes, opts...)
	return RunConsumer(resource.Source(), []itemConsumer{clean})
}

func InsertTask(source job.Table, target job.Table, mode transform.InsertOperator, opts ...func(*Insert)) error {
	insert := NewInsert(target, mode, opts...)
	return RunProducer(source, []itemProducer{insert})
}

func TransferTask(source job.Table, target job.Table, mode transform.TransferOperator) error {
	transfer := NewTransfer(target, mode)
	return RunProducer(source, []itemProducer{transfer})
}

func UpdateTransferTask(source job.Table, target job.Table, modes []transform.UpdateOperator) error {
	updateTransfer := NewUpdateTransfer(source, target, modes)
	return RunProducer(source, []itemProducer{updateTransfer})
}

func UpdateTask(source job.Table, modes []transform.UpdateOperator) error {
	update := NewUpdate(source, modes)
	return RunProducer(source, []itemProducer{update})
}
