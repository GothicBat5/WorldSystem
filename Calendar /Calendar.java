import javax.swing.*;
import java.awt.*;
import java.time.LocalDate;
import java.time.YearMonth;

public class Calendar extends JPanel
{
    private final Logic logic;

    private final String[] DAYS = {
        "Sun","Mon","Tue","Wed",
        "Thu","Fri","Sat"
    };

    public Calendar(Logic logic)
    {
        this.logic = logic;

        setBackground(Color.WHITE);
    }

    @Override
    protected void paintComponent(Graphics g)
    {
        super.paintComponent(g);

        Graphics2D g2 = (Graphics2D) g;

        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);

        int width = getWidth();
        int height = getHeight();

        int cellWidth = width / 7;
        int cellHeight = height / 7;

        g2.setFont(new Font("SansSerif", Font.BOLD, 16));

        for (int i = 0; i < 7; i++)
        {
            g2.drawString(DAYS[i], i * cellWidth + 15,  30);
        }

        YearMonth month = logic.getCurrentMonth();

        LocalDate today = logic.getToday();

        int start = logic.getStartColumn();

        int days = logic.getDaysInMonth();

        int day = 1;

        for (int row = 1; row < 7; row++)
        {
            for (int col = 0; col < 7; col++)
            {
                int x = col * cellWidth;
                int y = row * cellHeight;

                g2.setColor(Color.LIGHT_GRAY);
                g2.drawRect(x, y, cellWidth, cellHeight);

                if (row == 1 && col < start)
                {
                    continue;
                }

                if (day > days)
                {
                    break;
                }

                boolean isToday = today.getYear() == month.getYear() &&
                        today.getMonthValue() == month.getMonthValue() &&
                        today.getDayOfMonth() == day;

                if (isToday)
                {
                    g2.setColor(new Color(120, 180, 255));
                    g2.fillRoundRect(x + 5, y + 5, cellWidth - 10, cellHeight - 10, 20, 20);
                    g2.setColor(Color.BLACK);
                }

                g2.drawString(String.valueOf(day), x + 10, y + 25);

                day++;
            }
        }
    }
}
