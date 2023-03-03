package mysql

import (
	. "github.com/sokool/stream"
	"gorm.io/gorm"
)

type Table[E Entity] struct {
	name   string
	create NewEntity[E]
	c      *Connection
}

func NewTable[E Entity](c *Connection, fn NewEntity[E]) (*Table[E], error) {
	var e E

	t, err := NewType(e)
	if err != nil {
		return nil, err
	}

	d := &Table[E]{
		c:      c,
		create: fn,
		name:   t.LowerCase().String(),
	}

	return d, d.prepare(e, false)
}

func (r *Table[E]) Create(e Events) (E, error) {
	return r.create(e)
}

func (r *Table[E]) One(e E) error {
	db := r.c.gdb.Table(r.name).Where("id = ?", e.ID()).First(e)
	if err := db.Error; err != gorm.ErrRecordNotFound && err != nil {
		return err
	}

	return nil
}

func (r *Table[E]) Read(ee []E, bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *Table[E]) Update(ee ...E) (err error) {
	if len(ee) == 0 {
		return nil
	}

	if len(ee) == 1 {
		return r.c.gdb.Table(r.name).Save(ee[0]).Error
	}

	var tx *gorm.DB
	if tx = r.c.gdb.Table(r.name).Begin(); tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for i := range ee {
		if err = tx.Save(ee[i]).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (r *Table[E]) Delete(ee ...E) error {
	//return r.c.gdb.Table(r.name).Model(r.model).Where("id = ?", id).Delete(c.model).Error
	panic("implement")
}

func (r *Table[E]) prepare(e E, drop bool) error {
	db := r.c.gdb.
		Set("CHARACTER", "utf8mb4,utf8").
		Set("collation", "utf8mb4_unicode_ci").
		Table(r.name)

	if drop {
		if err := db.Migrator().DropTable(e); err != nil {
			return err
		}
	}

	return db.AutoMigrate(e)
}
