import javax.swing.JEditorPane; // The html loader
import javax.swing.text.html.HTMLEditorKit;
import java.io.IOException;
import java.net.URL;

public class HTML 
{
    private JEditorPane editorPane;

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

    public void loadPage(String url) 
    {
        
        try {
            String fixedUrl = url;
            if (!fixedUrl.startsWith("http://") && !fixedUrl.startsWith("https://")) 
            {
                fixedUrl = "https://" + fixedUrl;
            }
            editorPane.setPage(new URL(fixedUrl));
        } 
        catch (IOException e) 
        {
            editorPane.setText("<html><body><h2>Failed to load page</h2><p>" + e.getMessage()+"</p></body></html>");
        }
    }
}
