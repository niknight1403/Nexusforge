package dev.nexusforge.app;

import android.app.Activity;
import android.app.AlertDialog;
import android.content.Context;
import android.graphics.Bitmap;
import android.net.Uri;
import android.os.Bundle;
import android.view.ViewGroup;
import android.webkit.CookieManager;
import android.webkit.WebResourceError;
import android.webkit.WebResourceRequest;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import android.widget.Button;
import android.widget.EditText;
import android.widget.LinearLayout;
import android.widget.TextView;

public class MainActivity extends Activity {
    private WebView web;
    private TextView status;
    private String server;

    @Override
    public void onCreate(Bundle state) {
        super.onCreate(state);
        server = getPreferences(Context.MODE_PRIVATE).getString("server", "");
        if (server.isEmpty()) configure();
        else launch();
    }

    private static boolean validServer(String value) {
        Uri uri = Uri.parse(value);
        return "https".equalsIgnoreCase(uri.getScheme())
                && uri.getHost() != null && !uri.getHost().isEmpty()
                && uri.getUserInfo() == null
                && (uri.getPath() == null || uri.getPath().isEmpty() || "/".equals(uri.getPath()))
                && uri.getQuery() == null && uri.getFragment() == null
                && (uri.getPort() == -1 || (uri.getPort() > 0 && uri.getPort() <= 65535));
    }

    private void configure() {
        EditText input = new EditText(this);
        input.setSingleLine(true);
        input.setHint("https://nexusforge.example.com");
        input.setText(server);
        new AlertDialog.Builder(this)
                .setTitle("NexusForge Server verbinden")
                .setMessage("Gib die HTTPS-Adresse deines NexusForge-Servers ein.")
                .setView(input)
                .setCancelable(false)
                .setPositiveButton("Verbinden", (dialog, which) -> {
                    String value = input.getText().toString().trim().replaceAll("/+$", "");
                    if (!validServer(value)) {
                        new AlertDialog.Builder(this)
                                .setMessage("Bitte eine HTTPS-Serveradresse ohne Pfad, Zugangsdaten oder Parameter eingeben.")
                                .setPositiveButton("OK", (d, w) -> configure()).show();
                        return;
                    }
                    server = value;
                    getPreferences(Context.MODE_PRIVATE).edit().putString("server", server).apply();
                    launch();
                }).show();
    }

    private boolean isTrusted(Uri uri) {
        Uri trusted = Uri.parse(server);
        return "https".equalsIgnoreCase(uri.getScheme())
                && uri.getHost() != null
                && uri.getHost().equalsIgnoreCase(trusted.getHost())
                && uri.getPort() == trusted.getPort()
                && uri.getUserInfo() == null;
    }

    private void launch() {
        if (web != null) { web.destroy(); web = null; }
        LinearLayout layout = new LinearLayout(this);
        layout.setOrientation(LinearLayout.VERTICAL);

        LinearLayout toolbar = new LinearLayout(this);
        Button settings = new Button(this);
        settings.setText("Server wechseln");
        settings.setOnClickListener(v -> configure());
        toolbar.addView(settings, new LinearLayout.LayoutParams(0, ViewGroup.LayoutParams.WRAP_CONTENT, 1));
        Button retry = new Button(this);
        retry.setText("Neu laden");
        retry.setOnClickListener(v -> { if (web != null) web.reload(); });
        toolbar.addView(retry, new LinearLayout.LayoutParams(0, ViewGroup.LayoutParams.WRAP_CONTENT, 1));
        layout.addView(toolbar);

        status = new TextView(this);
        status.setPadding(16, 8, 16, 8);
        status.setText("Verbinde mit NexusForge …");
        layout.addView(status);

        web = new WebView(this);
        WebSettings options = web.getSettings();
        options.setJavaScriptEnabled(true);
        options.setDomStorageEnabled(true);
        options.setAllowFileAccess(false);
        options.setAllowContentAccess(false);
        options.setAllowFileAccessFromFileURLs(false);
        options.setAllowUniversalAccessFromFileURLs(false);
        options.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
        CookieManager.getInstance().setAcceptThirdPartyCookies(web, false);
        web.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                if (isTrusted(request.getUrl())) return false;
                status.setText("Externer Link aus Sicherheitsgründen blockiert.");
                return true;
            }

            @Override
            public void onPageStarted(WebView view, String url, Bitmap favicon) {
                status.setText("Lade " + Uri.parse(server).getHost() + " …");
            }

            @Override
            public void onPageFinished(WebView view, String url) {
                if (isTrusted(Uri.parse(url))) status.setText("Verbunden: " + Uri.parse(server).getHost());
            }

            @Override
            public void onReceivedError(WebView view, WebResourceRequest request, WebResourceError error) {
                if (request.isForMainFrame()) {
                    status.setText("Server nicht erreichbar. HTTPS-Adresse prüfen oder neu laden.");
                }
            }
        });
        layout.addView(web, new LinearLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, 0, 1));
        setContentView(layout);
        web.loadUrl(server);
    }

    @Override
    public void onBackPressed() {
        if (web != null && web.canGoBack()) web.goBack();
        else super.onBackPressed();
    }

    @Override
    protected void onDestroy() {
        if (web != null) { web.destroy(); web = null; }
        super.onDestroy();
    }
}
