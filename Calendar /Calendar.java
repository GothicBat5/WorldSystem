import javax.swing.*;
import java.awt.*;
import java.time.LocalDate;
import java.time.YearMonth;

public class Calendar extends JPanel
{
    private static final Color BG_PANEL = new Color(0xF8F9FA);
    private static final Color BG_HEADER = new Color(0xFFFFFF);
    private static final Color BG_CELL = new Color(0xFFFFFF);
    private static final Color BG_WEEKEND = new Color(0xFAFAFC);
    private static final Color BG_TODAY = new Color(0x4F6EF7);
    private static final Color BG_TODAY_CELL = new Color(0xEEF1FE);
    private static final Color LINE = new Color(0xE5E7EB);
    private static final Color TEXT_HEADER = new Color(0x9CA3AF);
    private static final Color TEXT_NORMAL = new Color(0x374151);
    private static final Color TEXT_MUTED = new Color(0xD1D5DB);
    private static final Color TEXT_TODAY = Color.WHITE;

    private static final String[] DAYS = { "Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat" };

    private final Logic logic;

    public Calendar(Logic logic)
    {
        this.logic = logic;
        setBackground(BG_PANEL);
    }

    @Override
    protected void paintComponent(Graphics g)
    {
        super.paintComponent(g);
        Graphics2D g2 = (Graphics2D) g;
        g2.setRenderingHint(RenderingHints.KEY_ANTIALIASING, RenderingHints.VALUE_ANTIALIAS_ON);
        g2.setRenderingHint(RenderingHints.KEY_TEXT_ANTIALIASING, RenderingHints.VALUE_TEXT_ANTIALIAS_LCD_HRGB);

        int w = getWidth();
        int h = getHeight();
        int headerH  = 40;
        int cellW = w / 7;
        int gridH = h - headerH;
        int cellH = gridH / 6;

        drawDayHeaders(g2, cellW, headerH);
        drawGrid(g2, cellW, cellH, headerH);
        drawDayNumbers(g2, cellW, cellH, headerH);
    }

    private void drawDayHeaders(Graphics2D g2, int cellW, int headerH)
    {
        g2.setColor(BG_HEADER);
        g2.fillRect(0, 0, getWidth(), headerH);

        g2.setColor(LINE);
        g2.drawLine(0, headerH, getWidth(), headerH);

        g2.setFont(new Font("Inter", Font.BOLD, 11));
        FontMetrics fm = g2.getFontMetrics();

        for (int i = 0; i < 7; i++)
        {
            boolean isWeekend = (i == 0 || i == 6);
            g2.setColor(isWeekend ? TEXT_MUTED : TEXT_HEADER);
            String label = DAYS[i];
            int textX = i * cellW + (cellW - fm.stringWidth(label)) / 2;
            int textY = (headerH + fm.getAscent()) / 2;
            g2.drawString(label, textX, textY);
        }
    }

    private void drawGrid(Graphics2D g2, int cellW, int cellH, int headerH)
    {
        for (int row = 0; row < 6; row++)
        {
            for (int col = 0; col < 7; col++)
            {
                int x = col * cellW;
                int y = headerH + row * cellH;
                boolean isWeekend = (col == 0 || col == 6);
                g2.setColor(isWeekend ? BG_WEEKEND : BG_CELL);
                g2.fillRect(x, y, cellW, cellH);
            }
        }
        // Vertical separators
        g2.setColor(LINE);
        for (int col = 1; col < 7; col++)
            g2.drawLine(col * cellW, headerH, col * cellW, getHeight());

        // Horizontal separators
        for (int row = 1; row < 6; row++)
            g2.drawLine(0, headerH + row * cellH, getWidth(), headerH + row * cellH);
    }

    private void drawDayNumbers(Graphics2D g2, int cellW, int cellH, int headerH)
    {
        YearMonth month = logic.getCurrentMonth();
        LocalDate today = logic.getToday();
        int start = logic.getStartColumn();
        int days = logic.getDaysInMonth();
        int day = 1;

        Font normalFont = new Font("Inter", Font.PLAIN, 13);
        Font todayFont  = new Font("Inter", Font.BOLD,  13);

        for (int row = 0; row < 6; row++)
        {
            for (int col = 0; col < 7; col++)
            {
                if (row == 0 && col < start) continue;
                if (day > days) return;

                int x = col * cellW;
                int y = headerH + row * cellH;

                boolean isToday = today.getYear() == month.getYear()
                               && today.getMonthValue() == month.getMonthValue()
                               && today.getDayOfMonth() == day;

                if (isToday)
                {
                    // Tint the whole cell lightly
                    g2.setColor(BG_TODAY_CELL);
                    g2.fillRect(x, y, cellW, cellH);

                    // Draw filled circle badge behind the number
                    int badgeSize = 30;
                    int bx = x + (cellW - badgeSize) / 2;
                    int by = y + 8;
                    g2.setColor(BG_TODAY);
                    g2.fillOval(bx, by, badgeSize, badgeSize);

                    g2.setFont(todayFont);
                    FontMetrics fm = g2.getFontMetrics();
                    String num = String.valueOf(day);
                    int tx = bx + (badgeSize - fm.stringWidth(num)) / 2;
                    int ty = by + (badgeSize + fm.getAscent() - fm.getDescent()) / 2;
                    g2.setColor(TEXT_TODAY);
                    g2.drawString(num, tx, ty);
                }
                else
                {
                    boolean isWeekend = (col == 0 || col == 6);
                    g2.setFont(normalFont);
                    FontMetrics fm = g2.getFontMetrics();
                    String num = String.valueOf(day);
                    int badgeSize = 28;
                    int bx = x + (cellW - badgeSize) / 2;
                    int by = y + 9;
                    int tx = bx + (badgeSize - fm.stringWidth(num)) / 2;
                    int ty = by + (badgeSize + fm.getAscent() - fm.getDescent()) / 2;
                    g2.setColor(isWeekend ? TEXT_MUTED : TEXT_NORMAL);
                    g2.drawString(num, tx, ty);
                }

                day++;
            }
        }
    }
}
