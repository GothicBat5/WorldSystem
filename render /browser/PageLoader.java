
public class PageLoader 
{
    private HTML htmlView;

    public PageLoader(HTML htmlView) 
    {
        this.htmlView = htmlView;
    }

    public void load(String url) 
    {
        htmlView.loadPage(url);
    }
}
