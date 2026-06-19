import javax.swing.JEditorPane; // The html loader reader ?? 
import javax.swing.SwingWorker;
import javax.swing.text.html.HTMLEditorKit;
import java.net.URL;

public class HTML 
{
    public interface PageLoadListener 
    {
        void onLoadStart(String url);
        void onLoadFinish(boolean success, String message);
    }

    private JEditorPane editorPane;
    private PageLoadListener listener;

    public HTML() 
    {
        editorPane = new JEditorPane();
        editorPane.setEditable(false);
        editorPane.setContentType("text/html");
        editorPane.setEditorKit(new HTMLEditorKit());
    }

    public JEditorPane getEditorPane() 
    {
        return editorPane;
    }

    public void setPageLoadListener(PageLoadListener listener) 
    {
        this.listener = listener;
    }

    public void loadPage(String url) 
    {
        String fixedUrl = url;
        if (!fixedUrl.startsWith("http://") && !fixedUrl.startsWith("https://")) 
        {
            fixedUrl = "https://" + fixedUrl;
        }
        final String finalUrl = fixedUrl;

        if (listener != null) 
        {
            listener.onLoadStart(finalUrl);
        }

        SwingWorker<Void, Void> worker = new SwingWorker<Void, Void>() {
            private Exception error;

            @Override
            protected Void doInBackground() 
            {
                try {
                    editorPane.setPage(new URL(finalUrl));
                } 
                catch (Exception e) {
                    error = e;
                }
                return null;
            }

            @Override
            protected void done() 
            {
                if (error != null) 
                {
                    
                    editorPane.setContentType("text/html");
                    editorPane.setText("<html><body><h2>Failed to load page</h2><p>" + error.getMessage()+"</p></body></html>");
                }
                if (listener != null) 
                {
                    listener.onLoadFinish(error == null, error == null ? finalUrl : error.getMessage());
                }
            }
        };

        worker.execute();
    }
}
