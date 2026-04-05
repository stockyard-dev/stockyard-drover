package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-drover/internal/server";"github.com/stockyard-dev/stockyard-drover/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="8720"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./drover-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("drover: %v",err)};defer db.Close();srv:=server.New(db,server.DefaultLimits())
fmt.Printf("\n  Drover — Self-hosted background job queue\n  ─────────────────────────────────\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n  Data:       %s\n  ─────────────────────────────────\n  Questions? hello@stockyard.dev\n\n",port,port,dataDir)
log.Printf("drover: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
