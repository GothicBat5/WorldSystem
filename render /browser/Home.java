import javax.swing.SwingUtilities;

public class Home 
{
    public static void main(String[] args) 
    {
        //Main Execution from :: Build
        //Try this: ( https://www.behance.net ) sample website from Adobe  
        // THis = import javax.swing.JEditorPane; The HTML renderer
        SwingUtilities.invokeLater(Build::new);
    }
}
