package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Queue struct{ID string `json:"id"`;Name string `json:"name"`;Description string `json:"description,omitempty"`;CreatedAt string `json:"created_at"`;PendingCount int `json:"pending_count"`;ProcessingCount int `json:"processing_count"`;DoneCount int `json:"done_count"`;FailedCount int `json:"failed_count"`}
type Job struct{ID string `json:"id"`;QueueID string `json:"queue_id"`;QueueName string `json:"queue_name,omitempty"`;Payload string `json:"payload,omitempty"`;Status string `json:"status"`;Priority int `json:"priority"`;Attempts int `json:"attempts"`;MaxAttempts int `json:"max_attempts"`;Error string `json:"error,omitempty"`;CreatedAt string `json:"created_at"`;StartedAt string `json:"started_at,omitempty"`;FinishedAt string `json:"finished_at,omitempty"`}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"drover.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
for _,q:=range[]string{
`CREATE TABLE IF NOT EXISTS queues(id TEXT PRIMARY KEY,name TEXT UNIQUE NOT NULL,description TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`,
`CREATE TABLE IF NOT EXISTS jobs(id TEXT PRIMARY KEY,queue_id TEXT NOT NULL,payload TEXT DEFAULT '',status TEXT DEFAULT 'pending',priority INTEGER DEFAULT 0,attempts INTEGER DEFAULT 0,max_attempts INTEGER DEFAULT 3,error TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')),started_at TEXT DEFAULT '',finished_at TEXT DEFAULT '')`,
`CREATE INDEX IF NOT EXISTS idx_jobs_queue ON jobs(queue_id)`,`CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status)`,
}{if _,err:=db.Exec(q);err!=nil{return nil,fmt.Errorf("migrate: %w",err)}};return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)CreateQueue(q *Queue)error{q.ID=genID();q.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO queues VALUES(?,?,?,?)`,q.ID,q.Name,q.Description,q.CreatedAt);return err}
func(d *DB)hydrateQueue(q *Queue){d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE queue_id=? AND status='pending'`,q.ID).Scan(&q.PendingCount);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE queue_id=? AND status='processing'`,q.ID).Scan(&q.ProcessingCount);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE queue_id=? AND status='done'`,q.ID).Scan(&q.DoneCount);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE queue_id=? AND status='failed'`,q.ID).Scan(&q.FailedCount)}
func(d *DB)GetQueue(id string)*Queue{var q Queue;if d.db.QueryRow(`SELECT id,name,description,created_at FROM queues WHERE id=?`,id).Scan(&q.ID,&q.Name,&q.Description,&q.CreatedAt)!=nil{return nil};d.hydrateQueue(&q);return &q}
func(d *DB)GetQueueByName(name string)*Queue{var q Queue;if d.db.QueryRow(`SELECT id,name,description,created_at FROM queues WHERE name=?`,name).Scan(&q.ID,&q.Name,&q.Description,&q.CreatedAt)!=nil{return nil};d.hydrateQueue(&q);return &q}
func(d *DB)ListQueues()[]Queue{rows,_:=d.db.Query(`SELECT id,name,description,created_at FROM queues ORDER BY name`);if rows==nil{return nil};defer rows.Close()
var o []Queue;for rows.Next(){var q Queue;rows.Scan(&q.ID,&q.Name,&q.Description,&q.CreatedAt);d.hydrateQueue(&q);o=append(o,q)};return o}
func(d *DB)DeleteQueue(id string)error{d.db.Exec(`DELETE FROM jobs WHERE queue_id=?`,id);_,err:=d.db.Exec(`DELETE FROM queues WHERE id=?`,id);return err}

func(d *DB)Enqueue(j *Job)error{j.ID=genID();j.CreatedAt=now();j.Status="pending";if j.MaxAttempts<=0{j.MaxAttempts=3}
_,err:=d.db.Exec(`INSERT INTO jobs(id,queue_id,payload,status,priority,attempts,max_attempts,created_at)VALUES(?,?,?,?,?,?,?,?)`,j.ID,j.QueueID,j.Payload,j.Status,j.Priority,0,j.MaxAttempts,j.CreatedAt);return err}
func(d *DB)Dequeue(queueID string)*Job{var j Job;if d.db.QueryRow(`SELECT id,queue_id,payload,status,priority,attempts,max_attempts,error,created_at,started_at,finished_at FROM jobs WHERE queue_id=? AND status='pending' ORDER BY priority DESC,created_at ASC LIMIT 1`,queueID).Scan(&j.ID,&j.QueueID,&j.Payload,&j.Status,&j.Priority,&j.Attempts,&j.MaxAttempts,&j.Error,&j.CreatedAt,&j.StartedAt,&j.FinishedAt)!=nil{return nil}
d.db.Exec(`UPDATE jobs SET status='processing',started_at=?,attempts=attempts+1 WHERE id=?`,now(),j.ID);j.Status="processing";j.StartedAt=now();return &j}
func(d *DB)Complete(id string)error{_,err:=d.db.Exec(`UPDATE jobs SET status='done',finished_at=? WHERE id=?`,now(),id);return err}
func(d *DB)Fail(id,errMsg string)error{var attempts,maxAttempts int;d.db.QueryRow(`SELECT attempts,max_attempts FROM jobs WHERE id=?`,id).Scan(&attempts,&maxAttempts)
status:="pending";if attempts>=maxAttempts{status="failed"}
_,err:=d.db.Exec(`UPDATE jobs SET status=?,error=?,finished_at=? WHERE id=?`,status,errMsg,now(),id);return err}
func(d *DB)Retry(id string)error{_,err:=d.db.Exec(`UPDATE jobs SET status='pending',error='',started_at='',finished_at='',attempts=0 WHERE id=?`,id);return err}
func(d *DB)ListJobs(queueID,status string,limit int)[]Job{if limit<=0{limit=50};q:=`SELECT j.id,j.queue_id,j.payload,j.status,j.priority,j.attempts,j.max_attempts,j.error,j.created_at,j.started_at,j.finished_at,COALESCE(q.name,'') FROM jobs j LEFT JOIN queues q ON j.queue_id=q.id WHERE j.queue_id=?`;args:=[]any{queueID}
if status!=""&&status!="all"{q+=` AND j.status=?`;args=append(args,status)};q+=` ORDER BY j.created_at DESC LIMIT ?`;args=append(args,limit)
rows,_:=d.db.Query(q,args...);if rows==nil{return nil};defer rows.Close()
var o []Job;for rows.Next(){var j Job;rows.Scan(&j.ID,&j.QueueID,&j.Payload,&j.Status,&j.Priority,&j.Attempts,&j.MaxAttempts,&j.Error,&j.CreatedAt,&j.StartedAt,&j.FinishedAt,&j.QueueName);o=append(o,j)};return o}
type Stats struct{Queues int `json:"queues"`;Pending int `json:"pending"`;Processing int `json:"processing"`;Done int `json:"done"`;Failed int `json:"failed"`}
func(d *DB)Stats()Stats{var s Stats;d.db.QueryRow(`SELECT COUNT(*) FROM queues`).Scan(&s.Queues);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE status='pending'`).Scan(&s.Pending);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE status='processing'`).Scan(&s.Processing);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE status='done'`).Scan(&s.Done);d.db.QueryRow(`SELECT COUNT(*) FROM jobs WHERE status='failed'`).Scan(&s.Failed);return s}
