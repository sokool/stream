package mysql

import (
	. "github.com/sokool/stream"
	"gorm.io/gorm"
	"reflect"
)

type Table[E Entity] struct {
	name   string
	create EntityFunc[E]
	c      *Connection
}

func NewTable[E Entity](c *Connection, fn EntityFunc[E]) (*Table[E], error) {
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
	d, err := r.create(e)
	if err != nil {
		return d, err
	}

	if reflect.ValueOf(d).IsNil() {
		return d, ErrDocumentNotSupported
	}

	if err = r.DB().Where("id = ?", d.ID()).First(&d).Error; err != gorm.ErrRecordNotFound && err != nil {
		return d, err
	}

	return d, nil
}

func (r *Table[E]) Read(query []byte) (ee []E, _ error) {
	return ee, r.DB().Raw(string(query)).Scan(&ee).Error
}

func (r *Table[E]) Update(ee ...E) (err error) {
	if len(ee) == 0 {
		return nil
	}

	if len(ee) == 1 {
		return r.DB().Save(ee[0]).Error
	}

	var tx *gorm.DB
	if tx = r.DB().Begin(); tx.Error != nil {
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
	db := r.DB().Set("CHARACTER", "utf8mb4,utf8").Set("collation", "utf8mb4_unicode_ci")
	if drop {
		if err := db.Migrator().DropTable(e); err != nil {
			return err
		}
	}

	return db.AutoMigrate(e)
}

func (r *Table[E]) DB() *gorm.DB {
	return r.c.gdb.Table(r.name)
}
