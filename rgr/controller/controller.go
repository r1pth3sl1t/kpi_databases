package controller

import (
	"fmt"
	"rgr/model"
	"rgr/view"
)

type Controller struct {
	view    *view.View
	model   *model.Model
	options []func()
}

func New() (*Controller, error) {
	c := new(Controller)
	var err error = nil
	c.model, err = model.New()
	if err != nil {
		return nil, err
	}
	c.view = view.New(c.model.FetchTableData(), c.model.FetchTablePrimaryKeys())
	c.options = make([]func(), 6)

	c.options[1] = c.InsertData
	c.options[2] = c.UpdateData
	c.options[3] = c.DeleteData
	c.options[4] = c.GenerateData
	c.options[5] = c.SearchData
	return c, err
}

func (c *Controller) Destroy() {
	fmt.Println("Closing connection")
	c.model.Close()
}

func (c *Controller) Index() bool {
	option, err := c.view.Index(len(c.options))
	if option == 0 {
		return false
	}
	if err != nil {
		c.view.Error(err)
		return true
	}
	c.options[option]()
	return true
}

func (c *Controller) InsertData() {
	table, err := c.view.SelectTable()
	if err != nil {
		c.view.Error(err)
		return
	}
	data := c.view.FetchAttributes(c.view.Tables[table])
	err = c.model.Insert(table, data)
	if err != nil {
		c.view.Error(err)
		return
	}
}

func (c *Controller) UpdateData() {
	table, err := c.view.SelectTable()
	if err != nil {
		c.view.Error(err)
		return
	}
	pkey := c.view.FetchPrimaryKey(table)
	columns := c.view.SelectAttributes(table)
	data := c.view.FetchAttributes(columns)

	err = c.model.Update(table, data, pkey)

	if err != nil {
		c.view.Error(err)
		return
	}
}

func (c *Controller) DeleteData() {
	table, err := c.view.SelectTable()
	if err != nil {
		c.view.Error(err)
	}

	pkey := c.view.FetchPrimaryKey(table)
	err = c.model.Delete(table, pkey)

	if err != nil {
		c.view.Error(err)
		return
	}
}

func (c *Controller) GenerateData() {
	num := c.view.GetDataSize()
	err := c.model.GenerateDataSet(num)
	if err != nil {
		c.view.Error(err)
		return
	}
}

func (c *Controller) SearchData() {
	searchHandler, err := c.view.GetSearchingMode()
	if err != nil {
		c.view.Error(err)
		return
	}
	data := c.view.FetchAttributes(searchHandler.FetchSearchAttributes())
	time, columns, err := c.model.Search(searchHandler.Search(), data, searchHandler.FetchSearchAttributes())

	if err != nil {
		c.view.Error(err)
		return
	}

	c.view.IndexColumns(time, columns, searchHandler.Head())

}
