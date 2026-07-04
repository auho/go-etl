package task

import (
	"github.com/auho/go-etl/v3/job"
	"github.com/auho/go-etl/v3/job/transform"
)

// CleanTask runs a clean (filter) task: rows matching the operators are moved to
// the deleted table, the rest stay in the data table.
func CleanTask(resource job.CleanResource, operators []transform.UpdateOperator, opts ...func(clean *Clean)) error {
	clean := NewClean(resource, operators, opts...)
	return RunConsumer(resource.Source(), []itemConsumer{clean})
}

// InsertTask runs an insert task: rows from source are transformed by operator and
// inserted into target.
func InsertTask(source job.Table, target job.Table, operator transform.InsertOperator, opts ...func(*Insert)) error {
	insert := NewInsert(target, operator, opts...)
	return RunProducer(source, []itemProducer{insert})
}

// TransferTask runs a transfer task: each source row is transformed by operator and
// written to target (target is truncated first).
func TransferTask(source job.Table, target job.Table, operator transform.TransferOperator) error {
	transfer := NewTransfer(target, operator)
	return RunProducer(source, []itemProducer{transfer})
}

// UpdateTransferTask runs an update+transfer task: source rows are updated in
// place by operators and the updated rows are copied to target.
func UpdateTransferTask(source job.Table, target job.Table, operators []transform.UpdateOperator) error {
	updateTransfer := NewUpdateTransfer(source, target, operators)
	return RunProducer(source, []itemProducer{updateTransfer})
}

// UpdateTask runs an in-place update task: source rows are updated by operators and
// written back to the source table.
func UpdateTask(source job.Table, operators []transform.UpdateOperator) error {
	update := NewUpdate(source, operators)
	return RunProducer(source, []itemProducer{update})
}