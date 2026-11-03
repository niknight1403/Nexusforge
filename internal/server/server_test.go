package server

import (
 "context"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
 "time"
)

func TestEndpoints(t *testing.T) {
 api:=&API{Token:"test-token", Router:Router{Providers:nil}}
 h:=api.Handler()
 for _,tc:=range []struct{path,method,token,body string;want int}{
  {"/healthz","GET","", "",200},
  {"/","GET","", "",200},
  {"/api/providers","GET","", "",401},
  {"/api/providers","GET","test-token","",200},
  {"/api/chat","POST","test-token",`{"messages":[]}`,400},
  {"/api/chat","POST","test-token",`{"messages":[{"role":"user","content":"hello"}]}`,503},
 }{
  req:=httptest.NewRequest(tc.method,tc.path,strings.NewReader(tc.body));req.Header.Set("X-NexusForge-Token",tc.token)
  w:=httptest.NewRecorder();h.ServeHTTP(w,req)
  if w.Code!=tc.want {t.Errorf("%s %s: got %d want %d: %s",tc.method,tc.path,w.Code,tc.want,w.Body.String())}
 }
}
func TestFallback(t *testing.T){
 bad:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){http.Error(w,"rate limited",429)}));defer bad.Close()
 good:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Header.Get("Authorization")!="Bearer secret" {t.Errorf("missing bearer token")}
  var v map[string]any;if err:=json.NewDecoder(r.Body).Decode(&v);err!=nil{t.Errorf("bad JSON: %v",err)}
  w.Header().Set("Content-Type","application/json");_,_=w.Write([]byte(`{"choices":[{"message":{"content":"ready"}}]}`))
 }));defer good.Close()
 router:=Router{Providers:[]Provider{{Name:"first",Endpoint:bad.URL,Model:"test",Key:"secret"},{Name:"second",Endpoint:good.URL,Model:"test2",Key:"secret"}},Client:&http.Client{Timeout:3*time.Second}}
 resp,err:=router.Chat(context.Background(),[]Message{{Role:"user",Content:"ping"}})
 if err!=nil || resp.Provider!="second" || resp.Content!="ready" {t.Fatalf("fallback: %+v %v",resp,err)}
}
