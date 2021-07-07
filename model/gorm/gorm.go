package gorm

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// postgres 驱动
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"

	"demo_casbin/config"
	"demo_casbin/model/gorm/entity"
)

const (
	uniqueShortIDFuncSQL = `CREATE EXTENSION IF NOT EXISTS "pgcrypto";
	CREATE OR REPLACE FUNCTION unique_short_id()
	RETURNS TRIGGER AS $$
	DECLARE
	  key TEXT;
	  qry TEXT;
	  found TEXT;
	BEGIN
	  qry := 'SELECT id FROM ' || quote_ident(TG_TABLE_NAME) || ' WHERE id=';
	  LOOP
		key := encode(gen_random_bytes(6), 'base64');
		key := replace(key, '/', 'x'); -- url safe replacement
		key := replace(key, '+', 'v'); -- url safe replacement
		EXECUTE qry || quote_literal(key) INTO found;
		IF found IS NULL THEN
		  EXIT;
		END IF;
	  END LOOP;
	  NEW.id = key;
	  RETURN NEW;
	END;
	$$ language 'plpgsql';`
)

func genShortIDTriggerSQL(tableName string) string {
	return fmt.Sprintf(`DO $do$ BEGIN IF EXISTS (SELECT 1 FROM pg_trigger WHERE  NOT tgisinternal AND tgname = 'trigger_%s_shortid') THEN ELSE
	CREATE TRIGGER trigger_%s_shortid BEFORE INSERT ON %s FOR EACH ROW EXECUTE PROCEDURE unique_short_id();
	END IF; END $do$`, tableName, tableName, tableName)
}

// NewDB 创建DB实例
func NewDB() (*gorm.DB, func(), error) {
	source := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		"127.0.0.1", 8412, "postgres", "casbin", "1", "disable")
	db, err := gorm.Open("postgres", source)
	if err != nil {
		return nil, nil, err
	}

	db = db.Debug()

	cleanFunc := func() {
		err := db.Close()
		if err != nil {
			logrus.Errorf("Gorm db close error: %s", err.Error())
		}
	}

	err = db.DB().Ping()
	if err != nil {
		return nil, cleanFunc, err
	}

	return db, cleanFunc, nil
}

// AutoMigrate 自动映射数据表
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		new(entity.User),
		new(entity.PolicySource),
		new(entity.Project),
		new(entity.ProjectUser),
	).Error
	if err != nil {
		return err
	}

	db.Exec(uniqueShortIDFuncSQL)
	db.Exec(genShortIDTriggerSQL("projects"))
	db.Exec(genShortIDTriggerSQL("policy_sources"))
	db.Exec(genShortIDTriggerSQL("users"))
	return nil
}

// InitGormDB 初始化gorm存储
func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C.Gorm
	db, cleanFunc, err := NewGormDB()
	if err != nil {
		return nil, cleanFunc, err
	}

	if cfg.EnableAutoMigrate {
		err = AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}
	}

	return db, cleanFunc, nil
}

// NewGormDB 创建DB实例
func NewGormDB() (*gorm.DB, func(), error) {
	//cfg := config.C
	//var dsn string
	//switch cfg.Gorm.DBType {
	//case "mysql":
	//	dsn = cfg.MySQL.DSN()
	//case "sqlite3":
	//	dsn = cfg.Sqlite3.DSN()
	//	_ = os.MkdirAll(filepath.Dir(dsn), 0777)
	//case "postgres":
	//	dsn = cfg.Postgres.DSN()
	//default:
	//	return nil, nil, errors.New("unknown db")
	//}

	return NewDB()
}
