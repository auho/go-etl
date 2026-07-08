package etl

import (
	"github.com/auho/go-etl/v3/task/transform"
)

// InsertTask runs an insert task: rows from source are transformed by operator and
// inserted into target.
func InsertTask(source Table, target Table, operator transform.InsertOperator, opts ...func(*Insert)) error {
	return runProducer(source, NewInsert(target, operator, opts...))
}

// CleanTask runs a clean (filter) task: rows matching the operators are moved to
// the deleted table, the rest stay in the data table.
func CleanTask(resource CleanResource, operators []transform.UpdateOperator, opts ...func(*Clean)) error {
	return runConsumer(resource.Source(), NewClean(resource, operators, opts...))
}

// TransferTask runs a transfer task: each source row is transformed by operator and
// written to target (target is truncated first).
func TransferTask(source Table, target Table, operator transform.TransferOperator) error {
	return runProducer(source, NewTransfer(target, operator))
}

// UpdateTransferTask runs an update+transfer task: source rows are updated in
// place by operators and the updated rows are copied to target.
func UpdateTransferTask(source Table, target Table, operators []transform.UpdateOperator, opts ...func(*UpdateTransfer)) error {
	return runProducer(source, NewUpdateTransfer(source, target, operators, opts...))
}

// UpdateTask runs an in-place update task: source rows are updated by operators and
// written back to the source table.
func UpdateTask(source Table, operators []transform.UpdateOperator, opts ...func(*Update)) error {
	return runProducer(source, NewUpdate(source, operators, opts...))
}
