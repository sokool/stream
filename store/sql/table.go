package sql

import (
	. "github.com/sokool/stream"
	"gorm.io/gorm"
)

type Table[E Entity] struct {
	*Connection
	name string
}

func NewTable[E Entity](c *Connection) (*Table[E], error) {
	var e E

	d := &Table[E]{
		c,
		MustType[E]().ToLower(),
	}

	return d, d.prepare(e, false)
}

func (r *Table[E]) One(e E) error {
	db := r.gdb.Table(r.name).Where("id = ?", e.ID()).First(e)
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
		return r.gdb.Table(r.name).Save(ee[0]).Error
	}

	var tx *gorm.DB
	if tx = r.gdb.Table(r.name).Begin(); tx.Error != nil {
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
	db := r.gdb.
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
