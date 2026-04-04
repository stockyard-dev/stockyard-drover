package server
import "net/http"
func(s *Server)dashboard(w http.ResponseWriter,r *http.Request){w.Header().Set("Content-Type","text/html");w.Write([]byte(dashHTML))}
const dashHTML=`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Drover</title>
<style>:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}
.main{padding:1.5rem;max-width:900px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1.2rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}.st-v{font-size:1.1rem}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.1rem}
.queue-bar{display:flex;gap:.3rem;margin-bottom:1rem}
.q-btn{font-size:.6rem;padding:.2rem .6rem;border:1px solid var(--bg3);background:var(--bg);color:var(--cm);cursor:pointer}.q-btn:hover{border-color:var(--leather)}.q-btn.active{border-color:var(--gold);color:var(--gold)}
.job{display:flex;justify-content:space-between;align-items:center;padding:.5rem .8rem;border-bottom:1px solid var(--bg3);font-size:.72rem}
.job:hover{background:var(--bg2)}
.status-pending{color:var(--gold)}.status-running{color:var(--orange)}.status-completed{color:var(--green)}.status-failed{color:var(--red)}.status-retrying{color:var(--orange)}
.badge{font-size:.5rem;padding:.1rem .3rem;text-transform:uppercase;letter-spacing:1px}
.badge-pending{background:#d4a84322;color:var(--gold);border:1px solid #d4a84344}
.badge-running{background:#d4843a22;color:var(--orange);border:1px solid #d4843a44}
.badge-completed{background:#4a9e5c22;color:var(--green);border:1px solid #4a9e5c44}
.badge-failed{background:#c9444422;color:var(--red);border:1px solid #c9444444}
.btn{font-size:.6rem;padding:.25rem .6rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd)}.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:var(--bg)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.6);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:400px;max-width:90vw}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust)}
.fr{margin-bottom:.5rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.15rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.35rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:.8rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1>DROVER</h1><div><button class="btn" onclick="openQueue()">+ Queue</button> <button class="btn btn-p" onclick="openJob()">+ Enqueue Job</button></div></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="queue-bar" id="queues"></div>
<div id="jobs"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)cm()"><div class="modal" id="mdl"></div></div>
<script>
const A='/api';let queues=[],jobs=[],curQ='';
async function load(){const[q,s]=await Promise.all([fetch(A+'/queues').then(r=>r.json()),fetch(A+'/stats').then(r=>r.json())]);
queues=q.queues||[];
const pending=s.pending||0,running=s.running||0,completed=s.completed||0,failed=s.failed||0;
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v" style="color:var(--gold)">'+pending+'</div><div class="st-l">Pending</div></div><div class="st"><div class="st-v" style="color:var(--orange)">'+running+'</div><div class="st-l">Running</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+completed+'</div><div class="st-l">Done</div></div><div class="st"><div class="st-v" style="color:var(--red)">'+failed+'</div><div class="st-l">Failed</div></div>';
let qh='<button class="q-btn'+(curQ===''?' active':'')+'" onclick="setQ(\'\')">All</button>';
queues.forEach(q=>{qh+='<button class="q-btn'+(curQ===q.id?' active':'')+'" onclick="setQ(\''+q.id+'\')">'+esc(q.name)+'</button>';});
document.getElementById('queues').innerHTML=qh;loadJobs();}
function setQ(id){curQ=id;loadJobs();}
async function loadJobs(){const url=curQ?A+'/queues/'+curQ+'/jobs':A+'/jobs';const r=await fetch(url).then(r=>r.json());jobs=r.jobs||[];renderJobs();}
function renderJobs(){if(!jobs.length){document.getElementById('jobs').innerHTML='<div class="empty">No jobs in queue.</div>';return;}
let h='';jobs.forEach(j=>{
h+='<div class="job"><div><span class="badge badge-'+j.status+'">'+j.status+'</span> <span style="margin-left:.3rem;color:var(--cream)">'+esc((j.payload||'').substring(0,80))+'</span>';
h+='<div style="font-size:.55rem;color:var(--cm);margin-top:.2rem">Attempt '+j.attempts+'/'+j.max_attempts;
if(j.started_at)h+=' · Started '+ft(j.started_at);
if(j.finished_at)h+=' · Done '+ft(j.finished_at);
if(j.error)h+=' · <span style="color:var(--red)">'+esc(j.error)+'</span>';
h+='</div></div><div style="display:flex;gap:.3rem">';
if(j.status==='failed')h+='<button class="btn" onclick="retry(\''+j.id+'\')">Retry</button>';
h+='<button class="btn" onclick="del(\''+j.id+'\')" style="font-size:.5rem;color:var(--cm)">✕</button></div></div>';});
document.getElementById('jobs').innerHTML=h;}
async function retry(id){await fetch(A+'/jobs/'+id+'/retry',{method:'POST'});load();}
async function del(id){await fetch(A+'/jobs/'+id,{method:'DELETE'});load();}
function openQueue(){document.getElementById('mdl').innerHTML='<h2>New Queue</h2><div class="fr"><label>Name</label><input id="f-n" placeholder="e.g. email-send"></div><div class="fr"><label>Description</label><input id="f-d"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="subQ()">Create</button></div>';document.getElementById('mbg').classList.add('open');}
async function subQ(){await fetch(A+'/queues',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:document.getElementById('f-n').value,description:document.getElementById('f-d').value})});cm();load();}
function openJob(){let opts=queues.map(q=>'<option value="'+q.id+'">'+esc(q.name)+'</option>').join('');
document.getElementById('mdl').innerHTML='<h2>Enqueue Job</h2><div class="fr"><label>Queue</label><select id="f-q">'+opts+'</select></div><div class="fr"><label>Payload (JSON)</label><textarea id="f-p" rows="3" placeholder=\'{"action":"send_email","to":"user@example.com"}\'></textarea></div><div class="fr"><label>Priority (0=normal)</label><input id="f-pr" type="number" value="0"></div><div class="acts"><button class="btn" onclick="cm()">Cancel</button><button class="btn btn-p" onclick="subJ()">Enqueue</button></div>';document.getElementById('mbg').classList.add('open');}
async function subJ(){await fetch(A+'/jobs',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({queue_id:document.getElementById('f-q').value,payload:document.getElementById('f-p').value,priority:parseInt(document.getElementById('f-pr').value)||0})});cm();load();}
function cm(){document.getElementById('mbg').classList.remove('open');}
function ft(t){if(!t)return'';return new Date(t).toLocaleTimeString([],{hour:'2-digit',minute:'2-digit'});}
function esc(s){if(!s)return'';const d=document.createElement('div');d.textContent=s;return d.innerHTML;}
load();setInterval(load,10000);
</script></body></html>`
