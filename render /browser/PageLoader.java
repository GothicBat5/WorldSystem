public class PageLoader 
{
    private HTML htmlView;

    public PageLoader(HTML htmlView) 
    {
        this.htmlView = htmlView;
    }

    public void load(String url) 
    {
        htmlView.loadPage(normalize(url));
    }

    private String normalize(String url) 
    {
        String trimmed = url.trim();
        String lower = trimmed.toLowerCase();
        
        if (lower.startsWith("http://") || lower.startsWith("https://") 
        || lower.startsWith("file://")) 
        {
            return trimmed;
        }
        return "https://" + trimmed;
    }
}
