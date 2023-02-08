// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"ariga.io/atlas/cmd/atlas/internal/migrate/ent/revision"
	"ariga.io/atlas/cmd/atlas/internal/migrate/ent/schema"
	"ariga.io/atlas/sql/migrate"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	revisionFields := schema.Revision{}.Fields()
	_ = revisionFields
	// revisionDescType is the schema descriptor for type field.
	revisionDescType := revisionFields[2].Descriptor()
	// revision.DefaultType holds the default value on creation for the type field.
	revision.DefaultType = migrate.RevisionType(revisionDescType.Default.(uint))
	// revisionDescApplied is the schema descriptor for applied field.
	revisionDescApplied := revisionFields[3].Descriptor()
	// revision.DefaultApplied holds the default value on creation for the applied field.
	revision.DefaultApplied = revisionDescApplied.Default.(int)
	// revision.AppliedValidator is a validator for the "applied" field. It is called by the builders before save.
	revision.AppliedValidator = revisionDescApplied.Validators[0].(func(int) error)
	// revisionDescTotal is the schema descriptor for total field.
	revisionDescTotal := revisionFields[4].Descriptor()
	// revision.DefaultTotal holds the default value on creation for the total field.
	revision.DefaultTotal = revisionDescTotal.Default.(int)
	// revision.TotalValidator is a validator for the "total" field. It is called by the builders before save.
	revision.TotalValidator = revisionDescTotal.Validators[0].(func(int) error)
}