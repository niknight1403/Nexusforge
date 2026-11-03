package server

import (
 "bytes"
 "context"
 "crypto/subtle"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "os"
 "strings"
 "time"
)

type Provider struct { Name string `json:"name"`; Model string `json:"model"`; Endpoint string `json:"-"`; Key string `json:"-"` }
type Message struct { Role string `json:"role"`; Content string `json:"content"` }
type ChatRequest struct { Messages []Message `json:"messages"` }
type ChatResponse struct { Provider string `json:"provider"`; Model string `json:"model"`; Content string `json:"content"` }
type Router struct { Providers []Provider; Client *http.Client }
type API struct { Router Router; Token string }
var errUnavailable = errors.New("all configured routes unavailable")

func configured() []Provider {
 options := []struct{name, model, endpoint, key string}{
  {"Groq", env("GROQ_MODEL", "llama-3.3-70b-versatile"), "https://api.groq.com/openai/v1/chat/completions", os.Getenv("GROQ_API_KEY")},
  {"OpenRouter", env("OPENROUTER_MODEL", "meta-llama/llama-3.3-70b-instruct:free"), "https://openrouter.ai/api/v1/chat/completions", os.Getenv("OPENROUTER_API_KEY")},
  {"Gemini", env("GEMINI_MODEL", "gemini-2.0-flash"), "https://generativelanguage.googleapis.com/v1beta/openai/chat/completions", os.Getenv("GEMINI_API_KEY")},
 }
 out := make([]Provider, 0, len(options))
 for _, v := range options { if v.key != "" { out = append(out, Provider{Name:v.name, Model:v.model, Endpoint:v.endpoint, Key:v.key}) } }
 return out
}
func env(k, fallback string) string { if v := os.Getenv(k); v != "" { return v }; return fallback }
func New() *API { return &API{Router:Router{Providers:configured(), Client:&http.Client{Timeout:25*time.Second}}, Token:os.Getenv("NEXUSFORGE_ACCESS_TOKEN")} }

func (a *API) Handler() http.Handler {
 mux := http.NewServeMux()
 mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request){ reply(w, 200, map[string]any{"status":"ok"}) })
 mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request){ if r.URL.Path != "/" { http.NotFound(w,r);return }; w.Header().Set("Content-Type","text/html; charset=utf-8"); w.Header().Set("Content-Security-Policy","default-src 'none'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'"); io.WriteString(w, indexHTML) })
 mux.HandleFunc("GET /api/providers", a.protected(func(w http.ResponseWriter, r *http.Request){ reply(w, 200, map[string]any{"providers":a.Router.Names()}) }))
 mux.HandleFunc("POST /api/chat", a.protected(a.chat))
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
  w.Header().Set("X-Content-Type-Options", "nosniff")
  w.Header().Set("Referrer-Policy","no-referrer")
  mux.ServeHTTP(w,r)
 })
}
func (a *API) protected(next http.HandlerFunc) http.HandlerFunc {
 return func(w http.ResponseWriter, r *http.Request) {
  if a.Token != "" && subtle.ConstantTimeCompare([]byte(r.Header.Get("X-NexusForge-Token")), []byte(a.Token)) != 1 { reply(w,401,map[string]string{"error":"invalid access token"});return }
  next(w,r)
 }
}
func (a *API) chat(w http.ResponseWriter, r *http.Request) {
 r.Body = http.MaxBytesReader(w,r.Body,64<<10)
 defer r.Body.Close()
 var req ChatRequest
 if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil || len(req.Messages)==0 || len(req.Messages)>30 { reply(w,400,map[string]string{"error":"invalid chat request"});return }
 for _,m:=range req.Messages { if (m.Role!="user"&&m.Role!="assistant"&&m.Role!="system") || strings.TrimSpace(m.Content)=="" || len(m.Content)>16000 { reply(w,400,map[string]string{"error":"invalid message"});return } }
 ctx,cancel:=context.WithTimeout(r.Context(),55*time.Second);defer cancel()
 result,err:=a.Router.Chat(ctx,req.Messages)
 if err!=nil { reply(w,503,map[string]string{"error":"no configured provider returned a response"});return }
 reply(w,200,result)
}
func (r Router) Names() []map[string]string {
 result:=make([]map[string]string,0,len(r.Providers))
 for _,p:=range r.Providers { result=append(result,map[string]string{"name":p.Name,"model":p.Model}) }
 return result
}
func (r Router) Chat(ctx context.Context, messages []Message) (ChatResponse,error) {
 client:=r.Client; if client==nil {client=&http.Client{Timeout:25*time.Second}}
 for _,p:=range r.Providers {
  reqbody,_:=json.Marshal(map[string]any{"model":p.Model,"messages":messages,"max_tokens":1024})
  req,err:=http.NewRequestWithContext(ctx,http.MethodPost,p.Endpoint,bytes.NewReader(reqbody));if err!=nil {continue}
  req.Header.Set("Authorization","Bearer "+p.Key)
  req.Header.Set("Content-Type","application/json")
  resp,err:=client.Do(req);if err!=nil {continue}
  limited:=io.LimitReader(resp.Body,1<<20); data,readErr:=io.ReadAll(limited);resp.Body.Close()
  if readErr!=nil || resp.StatusCode!=200 {continue}
  var parsed struct { Choices []struct { Message struct { Content string `json:"content"` } `json:"message"` } `json:"choices"` }
  if json.Unmarshal(data,&parsed)!=nil || len(parsed.Choices)==0 || parsed.Choices[0].Message.Content=="" {continue}
  return ChatResponse{Provider:p.Name,Model:p.Model,Content:parsed.Choices[0].Message.Content},nil
 }
 return ChatResponse{},fmt.Errorf("%w (%d configured)",errUnavailable,len(r.Providers))
}
func reply(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func Listen(addr string) error {
 a:=New()
 srv:=&http.Server{Addr:addr,Handler:a.Handler(),ReadHeaderTimeout:10*time.Second,IdleTimeout:60*time.Second}
 return srv.ListenAndServe()
}
