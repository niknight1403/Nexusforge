package server

import (
 "context"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestOllamaLocalConfiguration(t *testing.T) {
 t.Setenv("GROQ_API_KEY","")
 t.Setenv("OPENROUTER_API_KEY","")
 t.Setenv("GEMINI_API_KEY","")
 t.Setenv("OLLAMA_BASE_URL","http://127.0.0.1:11434")
 t.Setenv("OLLAMA_MODEL","qwen2.5:3b")
 got:=configured()
 if len(got)!=1 || got[0].Name!="Ollama" || got[0].Key!="" || got[0].Model!="qwen2.5:3b" || got[0].Endpoint!="http://127.0.0.1:11434/v1/chat/completions" {t.Fatalf("invalid Ollama configuration: %+v",got)}
}
func TestOllamaRejectsInsecureRemoteEndpoint(t *testing.T) {
 t.Setenv("GROQ_API_KEY","");t.Setenv("OPENROUTER_API_KEY","");t.Setenv("GEMINI_API_KEY","")
 t.Setenv("OLLAMA_BASE_URL","http://example.org:11434")
 if got:=configured();len(got)!=0 {t.Fatalf("insecure remote endpoint unexpectedly enabled: %+v",got)}
}
func TestUnauthenticatedRemoteBindRejected(t *testing.T) {
 t.Setenv("NEXUSFORGE_ACCESS_TOKEN","")
 if err:=Listen("0.0.0.0:8099");err==nil || !strings.Contains(err.Error(),"mandatory") {t.Fatalf("expected missing-token failure, got %v",err)}
}
func TestOllamaChatWithoutBearerToken(t *testing.T) {
 upstream:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if r.Header.Get("Authorization")!="" {t.Error("unexpected bearer token")}
  w.Header().Set("Content-Type","application/json")
  _,_ = w.Write([]byte(`{"choices":[{"message":{"content":"local works"}}]}`))
 }))
 defer upstream.Close()
 router:=Router{Providers:[]Provider{{Name:"Ollama",Model:"test",Endpoint:upstream.URL}}}
 res,err:=router.Chat(context.Background(),[]Message{{Role:"user",Content:"hi"}})
 if err!=nil || res.Content!="local works" {t.Fatalf("unexpected response: %+v %v",res,err)}
}
