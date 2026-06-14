import javax.swing.*;
import java.awt.*;

public class Frame extends JFrame
{
    private final Logic logic;
    private final Calendar calendar;
    private final JLabel titleLabel;

    public Frame()
    {
        logic = new Logic();

        calendar = new Calendar(logic);

        titleLabel = new JLabel("", SwingConstants.CENTER);

        JButton prev = new JButton("<");

        JButton next = new JButton(">");

        prev.addActionListener(e -> {
            logic.previousMonth(); 
            updateTitle();
            calendar.repaint();
        });

        next.addActionListener(e -> {
            logic.nextMonth();
            updateTitle();
            calendar.repaint();
        });

        JPanel top = new JPanel(new BorderLayout());

        top.add(prev, BorderLayout.WEST);
        top.add(titleLabel, BorderLayout.CENTER);
        top.add(next, BorderLayout.EAST);

        setLayout(new BorderLayout());

        add(top, BorderLayout.NORTH);
        add(calendar, BorderLayout.CENTER);

        updateTitle();

        setTitle("Java Calendar");
        setSize(900, 700);
        setLocationRelativeTo(null);
        setDefaultCloseOperation(EXIT_ON_CLOSE);
        setVisible(true);
    }

    private void updateTitle()
    {
        titleLabel.setText(logic.getCurrentMonth().getMonth()+" "+logic.getCurrentMonth().getYear());
    }
}
