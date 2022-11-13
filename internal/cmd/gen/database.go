package gen

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/db"
	"github.com/codfrm/cago/internal/cmd/gen/utils"
	"github.com/spf13/cobra"
)

const entityTpl = `
package entity

type {EntityName} struct {
{EntityField}
}
`

const repositoryTpl = `
package repository

import (
	"context"

	"{PkgName}"
)

type I{Name} interface {
	Find(ctx context.Context, id int64) (*entity.{Name}, error)
	Create(ctx context.Context, {LowerName} *entity.{Name}) error
	Update(ctx context.Context, {LowerName} *entity.{Name}) error
	Delete(ctx context.Context, id int64) error
}

var default{Name} I{Name}

func {Name}() I{Name} {
	return default{Name}
}

func Register{Name}(i I{Name}) {
	default{Name} = i
}
`

const persistenceTpl = `
package persistence

import (
	"context"

	"github.com/codfrm/cago/database/db"
	"{PkgName}/internal/model/entity"
	"{PkgName}/internal/repository"
)

type {LowerName} struct {
}

func New{Name}() repository.I{Name} {
	return &{LowerName}{}
}

func (u *{LowerName}) Find(ctx context.Context, id int64) (*entity.{Name}, error) {
	ret := &entity.{Name}{ID: id}
	if err := db.Ctx(ctx).First(ret).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return ret, nil
}

func (u *{LowerName}) Save(ctx context.Context, {LowerName} *entity.{Name}) error {
	return db.Ctx(ctx).Save({LowerName}).Error
}

func (u *{LowerName}) Update(ctx context.Context, {LowerName} *entity.{Name}) error {
	return db.Ctx(ctx).Updates({LowerName}).Error
}

func (u *{LowerName}) Delete(ctx context.Context, id int64) error {
	return db.Ctx(ctx).Delete(&entity.{Name}{ID: id}).Error
}
`

type Column struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type Index struct {
	Table      string
	NonUnique  int    `gorm:"column:Non_unique"`
	KeyName    string `gorm:"column:Key_name"`
	SeqInIndex int    `gorm:"column:Seq_in_index"`
	ColumnName string `gorm:"column:Column_name"`
}

func (c *Cmd) genDB(cmd *cobra.Command, args []string) error {
	table := args[0]
	// 读取appName
	abs, _ := filepath.Abs(".")
	cfg, err := configs.NewConfig(path.Base(abs))
	if err != nil {
		return err
	}
	c.pkgPath, c.pkgName, err = utils.FindRootPkgName("./")
	if err != nil {
		return err
	}
	if err := db.Database(context.Background(), cfg); err != nil {
		return err
	}
	column := make([]Column, 0)
	if err := db.Default().
		Raw("describe " + db.Default().Config.NamingStrategy.TableName(table)).Scan(&column).Error; err != nil {
		return err
	}
	index := make([]Index, 0)
	if err := db.Default().
		Raw("show index from " + db.Default().Config.NamingStrategy.TableName(table)).Scan(&index).Error; err != nil {
		return err
	}
	if err := c.genEntity(table, column, index); err != nil {
		return err
	}
	// 生成仓库接口
	if err := c.genRepository(table); err != nil {
		return err
	}
	// 生成仓库实现
	return c.genPersistence(table)
}

// 获取当前包名
func (c *Cmd) getCurrentPkgName(dir string) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	absPkgPath, err := filepath.Abs(c.pkgPath)
	if err != nil {
		return "", err
	}
	return c.pkgName + strings.TrimPrefix(absPath, absPkgPath), nil
}

func (c *Cmd) genRepository(table string) error {
	filepath := "internal/repository/" + table + ".go"
	// 存在不创建
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	// 获取表的索引信息
	repository := repositoryTpl
	pkgName, err := c.getCurrentPkgName("./internal/model/entity")
	if err != nil {
		return err
	}
	repository = strings.ReplaceAll(repository, "{PkgName}", pkgName)
	repository = strings.ReplaceAll(repository, "{Name}", utils.ToCamel(table))
	repository = strings.ReplaceAll(repository, "{LowerName}", utils.LowerFirstChar(utils.ToCamel(table)))
	// 写文件
	if err := os.MkdirAll("internal/repository", 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath, []byte(repository), 0644)
}

func (c *Cmd) genPersistence(table string) error {
	filepath := "internal/repository/persistence/" + table + ".go"
	// 存在不创建
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	persistence := persistenceTpl
	pkgName, err := c.getCurrentPkgName(".")
	if err != nil {
		return err
	}
	persistence = strings.ReplaceAll(persistence, "{PkgName}", pkgName)
	persistence = strings.ReplaceAll(persistence, "{Name}", utils.ToCamel(table))
	persistence = strings.ReplaceAll(persistence, "{LowerName}", utils.LowerFirstChar(utils.ToCamel(table)))
	// 写文件
	if err := os.MkdirAll("internal/repository/persistence", 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath, []byte(persistence), 0644)
}

func (c *Cmd) genEntity(table string, column []Column, index []Index) error {
	filepath := "internal/model/entity/" + table + ".go"
	// 存在不创建
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	// 获取表的索引信息
	entity := entityTpl
	entity = strings.ReplaceAll(entity, "{EntityName}", utils.ToCamel(table))
	var entityField string
	for _, v := range column {
		entityField += "\t" + utils.ToCamel(v.Field) + " " + convSqlType(v.Type) + " `" + convSqlTag(v, index) + "`\n"
	}
	entity = strings.ReplaceAll(entity, "{EntityField}", strings.TrimRight(entityField, "\n"))
	// 写文件
	if err := os.MkdirAll("internal/model/entity", 0755); err != nil {
		return err
	}
	return os.WriteFile(filepath, []byte(entity), 0644)
}

func convSqlTag(column Column, index []Index) string {
	tags := make([]string, 0)
	tags = append(tags, "column:"+column.Field)
	tags = append(tags, "type:"+column.Type)
	if column.Default != "" {
		tags = append(tags, "default:"+column.Default)
	}
	if column.Null == "NO" {
		tags = append(tags, "not null")
	}
	if column.Key == "PRI" {
		tags = append(tags, "primary_key")
	} else if column.Key != "" {
		// 生成索引
		for _, v := range index {
			if v.ColumnName == column.Field {
				tag := "index:" + v.KeyName
				if v.NonUnique == 0 {
					tag += ",unique"
				}
				tags = append(tags, tag)
			}
		}
	}
	return "gorm:\"" + strings.Join(tags, ";") + "\""
}

func convSqlType(sqlType string) string {
	// 取括号前的类型
	if strings.Contains(sqlType, "(") {
		sqlType = sqlType[:strings.Index(sqlType, "(")]
	}
	switch sqlType {
	case "int", "mediumint", "bigint":
		return "int64"
	case "tinyint", "smallint", "bit":
		return "int32"
	case "float", "double":
		return "float64"
	case "varchar", "char", "text", "mediumtext", "longtext":
		return "string"
	case "datetime", "timestamp":
		return "time.Time"
	}
	return "interface{}"
}
