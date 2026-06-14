import javax.swing.*;
import java.awt.*;
import java.time.format.TextStyle;
import java.util.Locale;

public class Frame extends JFrame
{
    private static final Color BG = new Color(0xF8F9FA);
    private static final Color HEADER_BG = Color.WHITE;
    private static final Color TEXT_DARK = new Color(0x1F2937);
    private static final Color ACCENT = new Color(0x4F6EF7);

    private final Logic logic;
    private final Calendar calendar;
    private final JLabel monthLabel;
    private final JLabel yearLabel;

    public Frame()
    {
        logic = new Logic();
        calendar = new Calendar(logic);
        monthLabel = makeLabel("", 22, Font.BOLD, TEXT_DARK);
        yearLabel = makeLabel("", 22, Font.PLAIN, new Color(0x6B7280));

        JPanel titleRow = new JPanel(new FlowLayout(FlowLayout.LEFT, 6, 0));
        titleRow.setOpaque(false);
        titleRow.add(monthLabel);
        titleRow.add(yearLabel);

        JPanel navRow = new JPanel(new FlowLayout(FlowLayout.RIGHT, 6, 0));
        navRow.setOpaque(false);
        navRow.add(navButton("\u2039", e -> { logic.previousMonth(); refresh(); }));
        navRow.add(navButton("\u203A", e -> { logic.nextMonth(); refresh(); }));

        JPanel header = new JPanel(new BorderLayout());
        header.setBackground(HEADER_BG);
        header.setBorder(BorderFactory.createCompoundBorder(BorderFactory.createMatteBorder(0, 0, 1, 0, new Color(0xE5E7EB)),
            BorderFactory.createEmptyBorder(18, 24, 18, 24)
        ));
        header.add(titleRow, BorderLayout.WEST);
        header.add(navRow, BorderLayout.EAST);

        getContentPane().setBackground(BG);
        setLayout(new BorderLayout());
        add(header, BorderLayout.NORTH);
        add(calendar, BorderLayout.CENTER);

        refresh();
        setTitle("Calendar");
        setSize(840, 680);
        setMinimumSize(new Dimension(600, 500));
        setLocationRelativeTo(null);
        setDefaultCloseOperation(EXIT_ON_CLOSE);
        setVisible(true);
    }

    private void refresh()
    {
        var ym = logic.getCurrentMonth();
        monthLabel.setText(ym.getMonth().getDisplayName(TextStyle.FULL, Locale.ENGLISH));
        yearLabel.setText(String.valueOf(ym.getYear()));
        calendar.repaint();
    }

    private JLabel makeLabel(String text, int size, int style, Color color)
    {
        JLabel label = new JLabel(text);
        label.setFont(new Font("Inter", style, size));
        label.setForeground(color);
        return label;
    }

    // Pill-shaped nav button matching modern calendar UX
    private JButton navButton(String symbol, java.awt.event.ActionListener action)
    {
        JButton btn = new JButton(symbol)
        {
            @Override
            protected void paintComponent(Graphics g)
            {
                Graphics2D g2 = (Graphics2D) g.create();
                g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
                if (getModel().isRollover())
                    g2.setColor(new Color(0xE9ECFF));
                else
                    g2.setColor(new Color(0xF3F4F6));
                g2.fillRoundRect(0, 0, getWidth(), getHeight(), getHeight(), getHeight());
                g2.dispose();
                super.paintComponent(g);
            }
        };
        btn.setFont(new Font("SansSerif", Font.PLAIN, 18));
        btn.setForeground(ACCENT);
        btn.setPreferredSize(new Dimension(36, 36));
        btn.setContentAreaFilled(false);
        btn.setBorderPainted(false);
        btn.setFocusPainted(false);
        btn.setCursor(Cursor.getPredefinedCursor(Cursor.HAND_CURSOR));
        btn.addActionListener(action);
        return btn;
    }
}
