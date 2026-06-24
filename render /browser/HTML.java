import javax.swing.JEditorPane;
import javax.swing.SwingWorker;
import javax.swing.text.Document;
import javax.swing.text.html.HTMLEditorKit;
import java.net.URL;

public class HTML 
{
    public interface PageLoadListener 
    {
        void onLoadStart(String url);
        void onLoadFinish(boolean success, String message);
        void onTitleChanged(String title);
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
        URL parsedUrl;
        
        try {
            parsedUrl = new URL(url);
        } 
        catch (Exception e) 
        {
            
            if (listener != null) 
            {
                listener.onLoadStart(url);
                listener.onLoadFinish(false, "Invalid URL: " + url);
            }
            editorPane.setContentType("text/html");
            editorPane.setText("<html><body><h2>Invalid URL</h2><p>" +url+ "</p></body></html>");
            return;
        }

        if (listener != null) 
        {
            listener.onLoadStart(url);
        }

        SwingWorker<Void, Void> worker = new SwingWorker<Void, Void>() 
        {
            private Exception error;

            @Override
            protected Void doInBackground() 
            {
                try {
                    editorPane.setPage(parsedUrl);
                } 
                catch (Exception e) 
                {
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
                    editorPane.setText("<html><body><h2>Failed to load page</h2><p>" + error.getMessage() + "</p></body></html>");
                    
                    if (listener != null) 
                    {
                        listener.onLoadFinish(false, error.getMessage());
                    }
                    return;
                }

                Object titleProperty = editorPane.getDocument().getProperty(Document.TitleProperty);
                String title = titleProperty != null ? titleProperty.toString().trim() : "";
                if (listener != null) 
                {
                    listener.onTitleChanged(title.isEmpty() ? url : title);
                    listener.onLoadFinish(true, url);
                }
            }
        };

        worker.execute();
    }
}
