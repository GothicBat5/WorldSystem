import javax.swing.JFrame;
import javax.swing.JTextField;
import javax.swing.JButton;
import javax.swing.JPanel;
import javax.swing.JScrollPane;
import java.awt.BorderLayout;
import java.awt.event.ActionEvent;

public class Build 
{
    private JFrame frame;
    private JTextField urlField;
    private JButton goButton;
    private HTML htmlView;
    private PageLoader pageLoader;

    public Build() 
    {
        htmlView = new HTML();
        pageLoader = new PageLoader(htmlView);

        frame = new JFrame("Mini Java Browser");
        frame.setSize(900, 600);
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        frame.setLayout(new BorderLayout());

        urlField = new JTextField();
        goButton = new JButton("Go");

        JPanel topPanel = new JPanel(new BorderLayout());
        topPanel.add(urlField, BorderLayout.CENTER);
        topPanel.add(goButton, BorderLayout.EAST);

        JScrollPane scrollPane = new JScrollPane(htmlView.getEditorPane());

        frame.add(topPanel, BorderLayout.NORTH);
        frame.add(scrollPane, BorderLayout.CENTER);

        goButton.addActionListener(this::onGoClicked);
        urlField.addActionListener(this::onGoClicked);

        frame.setVisible(true);
    }

    private void onGoClicked(ActionEvent e) 
    {
        String url = urlField.getText().trim();
        
        if (!url.isEmpty()) 
        {
            pageLoader.load(url);
        }
    }
}
