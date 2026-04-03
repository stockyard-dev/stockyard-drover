package server
import ("encoding/json";"log";"net/http";"github.com/stockyard-dev/stockyard-drover/internal/store")
type Server struct{db *store.DB;mux *http.ServeMux;limits Limits}
func New(db *store.DB,limits Limits)*Server{s:=&Server{db:db,mux:http.NewServeMux(),limits:limits}
s.mux.HandleFunc("GET /api/queues",s.listQueues);s.mux.HandleFunc("POST /api/queues",s.createQueue);s.mux.HandleFunc("GET /api/queues/{id}",s.getQueue);s.mux.HandleFunc("DELETE /api/queues/{id}",s.deleteQueue)
s.mux.HandleFunc("POST /api/queues/{id}/enqueue",s.enqueue);s.mux.HandleFunc("POST /api/queues/{id}/dequeue",s.dequeue)
s.mux.HandleFunc("GET /api/queues/{id}/jobs",s.listJobs)
s.mux.HandleFunc("POST /api/jobs/{id}/complete",s.complete);s.mux.HandleFunc("POST /api/jobs/{id}/fail",s.fail);s.mux.HandleFunc("POST /api/jobs/{id}/retry",s.retry)
s.mux.HandleFunc("GET /api/stats",s.stats);s.mux.HandleFunc("GET /api/health",s.health)
s.mux.HandleFunc("GET /ui",s.dashboard);s.mux.HandleFunc("GET /ui/",s.dashboard);s.mux.HandleFunc("GET /",s.root);
s.mux.HandleFunc("GET /api/tier",func(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"tier":s.limits.Tier,"upgrade_url":"https://stockyard.dev/drover/"})})
return s}
func(s *Server)ServeHTTP(w http.ResponseWriter,r *http.Request){s.mux.ServeHTTP(w,r)}
func wj(w http.ResponseWriter,c int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(c);json.NewEncoder(w).Encode(v)}
func we(w http.ResponseWriter,c int,m string){wj(w,c,map[string]string{"error":m})}
func(s *Server)root(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/"{http.NotFound(w,r);return};http.Redirect(w,r,"/ui",302)}
func(s *Server)listQueues(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"queues":oe(s.db.ListQueues())})}
func(s *Server)createQueue(w http.ResponseWriter,r *http.Request){var q store.Queue;json.NewDecoder(r.Body).Decode(&q);if q.Name==""{we(w,400,"name required");return};s.db.CreateQueue(&q);wj(w,201,s.db.GetQueue(q.ID))}
func(s *Server)getQueue(w http.ResponseWriter,r *http.Request){q:=s.db.GetQueue(r.PathValue("id"));if q==nil{we(w,404,"not found");return};wj(w,200,q)}
func(s *Server)deleteQueue(w http.ResponseWriter,r *http.Request){s.db.DeleteQueue(r.PathValue("id"));wj(w,200,map[string]string{"deleted":"ok"})}
func(s *Server)enqueue(w http.ResponseWriter,r *http.Request){qid:=r.PathValue("id");if s.db.GetQueue(qid)==nil{we(w,404,"queue not found");return}
var j store.Job;json.NewDecoder(r.Body).Decode(&j);j.QueueID=qid;s.db.Enqueue(&j);wj(w,201,j)}
func(s *Server)dequeue(w http.ResponseWriter,r *http.Request){j:=s.db.Dequeue(r.PathValue("id"));if j==nil{wj(w,200,map[string]any{"job":nil});return};wj(w,200,j)}
func(s *Server)listJobs(w http.ResponseWriter,r *http.Request){wj(w,200,map[string]any{"jobs":oe(s.db.ListJobs(r.PathValue("id"),r.URL.Query().Get("status"),100))})}
func(s *Server)complete(w http.ResponseWriter,r *http.Request){s.db.Complete(r.PathValue("id"));wj(w,200,map[string]string{"completed":"ok"})}
func(s *Server)fail(w http.ResponseWriter,r *http.Request){var req struct{Error string `json:"error"`};json.NewDecoder(r.Body).Decode(&req);s.db.Fail(r.PathValue("id"),req.Error);wj(w,200,map[string]string{"failed":"ok"})}
func(s *Server)retry(w http.ResponseWriter,r *http.Request){s.db.Retry(r.PathValue("id"));wj(w,200,map[string]string{"retried":"ok"})}
func(s *Server)stats(w http.ResponseWriter,r *http.Request){wj(w,200,s.db.Stats())}
func(s *Server)health(w http.ResponseWriter,r *http.Request){st:=s.db.Stats();wj(w,200,map[string]any{"status":"ok","service":"drover","pending":st.Pending,"processing":st.Processing})}
func oe[T any](s []T)[]T{if s==nil{return[]T{}};return s}
func init(){log.SetFlags(log.LstdFlags|log.Lshortfile)}
