package mysql

type DBLog struct {
	DbName string
}

func NewDBLog(dbName string) *DBLog {
	return &DBLog{
		DbName: dbName,
	}
}

func (this *DBLog) Event(eventName string) {
	return
}

func (this *DBLog) EventKv(eventName string, kvs map[string]string) {
	return
}

func (this *DBLog) EventErr(eventName string, err error) error {
	return err
}

func (this *DBLog) EventErrKv(eventName string, err error, kvs map[string]string) error {
	return err
}

func (this *DBLog) Timing(eventName string, nanoseconds int64) {
	return
}

func (this *DBLog) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	t := float32(nanoseconds) / float32(1000000)
	return
}
