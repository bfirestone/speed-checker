// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/bfirestone/speed-checker/ent/predicate"
	"github.com/bfirestone/speed-checker/ent/speedtest"
)

// SpeedTestDelete is the builder for deleting a SpeedTest entity.
type SpeedTestDelete struct {
	config
	hooks    []Hook
	mutation *SpeedTestMutation
}

// Where appends a list predicates to the SpeedTestDelete builder.
func (std *SpeedTestDelete) Where(ps ...predicate.SpeedTest) *SpeedTestDelete {
	std.mutation.Where(ps...)
	return std
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (std *SpeedTestDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, std.sqlExec, std.mutation, std.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (std *SpeedTestDelete) ExecX(ctx context.Context) int {
	n, err := std.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (std *SpeedTestDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(speedtest.Table, sqlgraph.NewFieldSpec(speedtest.FieldID, field.TypeInt))
	if ps := std.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, std.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	std.mutation.done = true
	return affected, err
}

// SpeedTestDeleteOne is the builder for deleting a single SpeedTest entity.
type SpeedTestDeleteOne struct {
	std *SpeedTestDelete
}

// Where appends a list predicates to the SpeedTestDelete builder.
func (stdo *SpeedTestDeleteOne) Where(ps ...predicate.SpeedTest) *SpeedTestDeleteOne {
	stdo.std.mutation.Where(ps...)
	return stdo
}

// Exec executes the deletion query.
func (stdo *SpeedTestDeleteOne) Exec(ctx context.Context) error {
	n, err := stdo.std.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{speedtest.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (stdo *SpeedTestDeleteOne) ExecX(ctx context.Context) {
	if err := stdo.Exec(ctx); err != nil {
		panic(err)
	}
}
