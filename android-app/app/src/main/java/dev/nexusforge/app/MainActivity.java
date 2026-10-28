package dev.nexusforge.app;
import android.app.Activity;
import android.app.AlertDialog;
import android.os.Bundle;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.webkit.WebResourceRequest;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.Button;
import android.view.ViewGroup;
import android.content.Context;
import android.net.Uri;
public class MainActivity extends Activity {
 private WebView web; private String server;
 @Override public void onCreate(Bundle state) { super.onCreate(state); server=getPreferences(Context.MODE_PRIVATE).getString("server",""); if(server.isEmpty())configure();else launch(); }
 private void configure() {
  EditText input=new EditText(this);input.setSingleLine(true);input.setHint("https://nexusforge.example.com");input.setText(server);
  new AlertDialog.Builder(this).setTitle("NexusForge Server verbinden").setMessage("Gib die HTTPS-Adresse deines NexusForge-Servers ein.").setView(input).setCancelable(false).setPositiveButton("Verbinden",(dialog,which)->{
   String value=input.getText().toString().trim().replaceAll("/+$","");Uri uri=Uri.parse(value);
   if(!"https".equalsIgnoreCase(uri.getScheme())||uri.getHost()==null){new AlertDialog.Builder(this).setMessage("Nur gültige HTTPS-Serveradressen sind erlaubt.").setPositiveButton("OK",(d,w)->configure()).show();return;}
   server=value;getPreferences(Context.MODE_PRIVATE).edit().putString("server",server).apply();launch();
  }).show();
 }
 private void launch() {
  LinearLayout layout=new LinearLayout(this);layout.setOrientation(LinearLayout.VERTICAL);
  Button settings=new Button(this);settings.setText("Server wechseln");settings.setOnClickListener(v->configure());layout.addView(settings,new LinearLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT,ViewGroup.LayoutParams.WRAP_CONTENT));
  web=new WebView(this);web.getSettings().setJavaScriptEnabled(true);web.getSettings().setDomStorageEnabled(true);
  web.setWebViewClient(new WebViewClient(){@Override public boolean shouldOverrideUrlLoading(WebView view,WebResourceRequest request){Uri u=request.getUrl();Uri trusted=Uri.parse(server);return !"https".equalsIgnoreCase(u.getScheme())||!trusted.getHost().equalsIgnoreCase(u.getHost())||trusted.getPort()!=u.getPort();}});
  layout.addView(web,new LinearLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT,0,1));setContentView(layout);web.loadUrl(server);
 }
 @Override public void onBackPressed(){if(web!=null&&web.canGoBack())web.goBack();else super.onBackPressed();}
}
