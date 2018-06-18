package paginator

import (
	"fmt"

	"github.com/go-zero-boilerplate/databases/sql_statement"
)

type ListContainer interface {
	Count() int
	SlicePointer() interface{}
	Clear()
	AfterSliceLoaded()
}

type DBPaginator interface {
	HasMore() bool
	GetNextIndex() (int, error)
}

func NewDBPaginator(selectBuilder sql_statement.SelectBuilder, pageSize int, listContainer ListContainer) (DBPaginator, error) {
	iterator := &dbPaginator{
		selectBuilder: selectBuilder,
		currentOffset: 0,
		pageSize:      pageSize,
		totalCount:    0,
		index:         0,
		listContainer: listContainer,
	}
	if err := iterator.loadTotalCount(); err != nil {
		return nil, err
	}

	if err := iterator.loadCurrentPage(); err != nil {
		return nil, err
	}

	return iterator, nil
}

type dbPaginator struct {
	selectBuilder sql_statement.SelectBuilder
	totalCount    int
	pageSize      int

	currentOffset int
	index         int

	listContainer ListContainer
}

func (d *dbPaginator) HasMore() bool {
	if d.index < d.listContainer.Count() {
		return true
	}
	return d.currentOffset+d.listContainer.Count() < d.totalCount
}

func (d *dbPaginator) GetNextIndex() (int, error) {
	if d.index < d.listContainer.Count() {
		nextIndex := d.index
		d.index++
		return nextIndex, nil
	}

	err := d.loadNextPage()
	if err != nil {
		return -1, err
	}

	if d.listContainer.Count() == 0 {
		return -1, fmt.Errorf("Paginator does not have more entries")
	}

	return d.GetNextIndex()
}

func (d *dbPaginator) loadNextPage() error {
	d.currentOffset += d.pageSize

	return d.loadCurrentPage()
}

func (d *dbPaginator) loadCurrentPage() error {
	d.index = 0
	d.listContainer.Clear()

	err := d.selectBuilder.Clone().ClearMappings().
		MapSlice(d.listContainer.SlicePointer()).
		Limit(d.pageSize).
		Offset(d.currentOffset).
		Build().
		Execute()
	if err != nil {
		return err
	}
	d.listContainer.AfterSliceLoaded()

	return nil
}

func (d *dbPaginator) loadTotalCount() error {
	var tmpCnt int
	countStatement := d.selectBuilder.Clone().ClearMappings().ClearLimit().ClearOffset().MapCountStar(&tmpCnt).Build()
	err := countStatement.Execute()
	if err != nil {
		return err
	}

	d.totalCount = tmpCnt
	return nil
}
