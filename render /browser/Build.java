import javax.swing.JFrame;
import javax.swing.JTextField;
import javax.swing.JButton;
import javax.swing.JPanel;
import javax.swing.JLabel;
import javax.swing.JScrollPane;
import javax.swing.BorderFactory;
import javax.swing.SwingConstants;
import java.awt.BorderLayout;
import java.awt.Color;
import java.awt.Font;
import java.awt.Dimension;
import java.awt.event.ActionEvent;

public class Build implements HTML.PageLoadListener  
{
    private JFrame frame;
    private JTextField urlField;
    private JButton goButton;
    private JLabel statusLabel;
    private HTML htmlView;
    private PageLoader pageLoader;

    private static final Color TOP_BAR_COLOR = new Color(245, 246, 248);
    private static final Color STATUS_BAR_COLOR = new Color(250, 250, 250);
    private static final Color ACCENT_COLOR = new Color(66, 133, 244);
    private static final Font UI_FONT = new Font("SansSerif", Font.PLAIN, 14);

    public Build() {
        htmlView = new HTML();
        htmlView.setPageLoadListener(this);
        pageLoader = new PageLoader(htmlView);

        frame = new JFrame("Mini Java Browser");
        frame.setMinimumSize(new Dimension(700, 500));
        frame.setSize(1000, 650);
        frame.setLocationRelativeTo(null);
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        frame.setLayout(new BorderLayout());

        urlField = new JTextField();
        urlField.setFont(UI_FONT);
        urlField.setBorder(BorderFactory.createCompoundBorder(BorderFactory.createLineBorder(new Color(210, 210, 210)),
                BorderFactory.createEmptyBorder(6, 10, 6, 10)
        ));

        goButton = new JButton("Go");
        goButton.setFont(UI_FONT);
        goButton.setFocusPainted(false);
        goButton.setBackground(ACCENT_COLOR);
        goButton.setForeground(Color.WHITE);
        goButton.setBorder(BorderFactory.createEmptyBorder(6, 18, 6, 18));

        JPanel topPanel = new JPanel(new BorderLayout(10, 0));
        topPanel.setBackground(TOP_BAR_COLOR);
        topPanel.setBorder(BorderFactory.createEmptyBorder(12, 14, 12, 14));
        topPanel.add(urlField, BorderLayout.CENTER);
        topPanel.add(goButton, BorderLayout.EAST);

        statusLabel = new JLabel("Ready");
        statusLabel.setFont(new Font("SansSerif", Font.PLAIN, 12));
        statusLabel.setForeground(new Color(120, 120, 120));
        statusLabel.setHorizontalAlignment(SwingConstants.LEFT);

        JPanel statusPanel = new JPanel(new BorderLayout());
        statusPanel.setBackground(STATUS_BAR_COLOR);
        statusPanel.setBorder(BorderFactory.createCompoundBorder(
                BorderFactory.createMatteBorder(1, 0, 0, 0, new Color(225, 225, 225)),
                BorderFactory.createEmptyBorder(5, 14, 5, 14)
        ));
        statusPanel.add(statusLabel, BorderLayout.WEST);

        JScrollPane scrollPane = new JScrollPane(htmlView.getEditorPane());
        scrollPane.setBorder(BorderFactory.createEmptyBorder());

        frame.add(topPanel, BorderLayout.NORTH);
        frame.add(scrollPane, BorderLayout.CENTER);
        frame.add(statusPanel, BorderLayout.SOUTH);

        goButton.addActionListener(this::onGoClicked);
        urlField.addActionListener(this::onGoClicked);

        frame.setVisible(true);
    }

    private void onGoClicked(ActionEvent e) {
        
        String url = urlField.getText().trim();
        if (!url.isEmpty()) 
        {
            pageLoader.load(url);
        }
    }

    @Override
    public void onLoadStart(String url) 
    {
        goButton.setEnabled(false);
        statusLabel.setForeground(new Color(120, 120, 120));
        statusLabel.setText("Loading " + url + " ...");
    }

    @Override
    public void onLoadFinish(boolean success, String message) 
    {
        goButton.setEnabled(true);
        
        if (success) 
        {
            statusLabel.setForeground(new Color(60, 140, 60));
            statusLabel.setText("Loaded " + message);
        } 
        else {
            statusLabel.setForeground(new Color(190, 60, 60));
            statusLabel.setText("Failed: " + message);
        }
    }
}
