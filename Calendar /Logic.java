import java.time.DayOfWeek;
import java.time.LocalDate;
import java.time.YearMonth;

public class Logic
{
    private YearMonth currentMonth;

    public Logic()
    {
        currentMonth = YearMonth.now();
    }

    public void nextMonth()
    {
        currentMonth = currentMonth.plusMonths(1);
    }

    public void previousMonth()
    {
        currentMonth = currentMonth.minusMonths(1);
    }

    public YearMonth getCurrentMonth()
    {
        return currentMonth;
    }

    public int getDaysInMonth()
    {
        return currentMonth.lengthOfMonth();
    }

    public int getStartColumn()
    {
        DayOfWeek day = currentMonth.atDay(1).getDayOfWeek();
        return day.getValue() % 7;
    }

    public LocalDate getToday()
    {
        return LocalDate.now();
    }
}
